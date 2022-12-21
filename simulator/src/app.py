import os
import shutil
import json
from datetime import datetime, timedelta

from flask import Flask, jsonify, Response, request
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import text, desc
from werkzeug.exceptions import BadRequest

import constants
from config_parser import parse_config, get_source_code_dir
from data_ingestion import State, DataIngestion
from cluster import Cluster
from open_search_simulator import Simulator
from plotter import plot_data_points


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///datapoints.db'
app.app_context().push()
if os.path.exists('instance'):
    shutil.rmtree('instance')
db = SQLAlchemy(app)


# Database model to store the datapoints
class DataModel(db.Model):
    status = db.Column(db.String(200))
    cpu_usage_percent = db.Column(db.Float, default=0)
    memory_usage_percent = db.Column(db.Float, default=0)
    shards_count = db.Column(db.Integer, default=0)
    total_nodes_count = db.Column(db.Integer, default=0)
    date_created = db.Column(db.DateTime, default=datetime.now(), primary_key=True)


# Converts the duration in minutes to time object of "HH:MM" format
def convert_to_hh_mm(duration_in_m):
    time_h_m = '{:02d}:{:02d}'.format(*divmod(duration_in_m, 60))
    time_obj = datetime.strptime(time_h_m, '%H:%M')
    return time_obj


# Returns the violated count for a requested metric, threshold and duration,
# returns error if sufficient data points are not present.
@app.route('/stats/violated/<string:stat_name>/<int:duration>/<float:threshold>')
def violated_count(stat_name, duration, threshold):
    # calculate time to query for data 
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points
    time_obj = time_now - timedelta(minutes=duration)

    try:
        # Fetching the count of data points for given duration.
        data_point_count = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(
            DataModel.date_created > time_obj).filter(DataModel.date_created < time_now).count()

        # If expected data points are not present then respond with error
        if duration // sim.frequency_minutes > data_point_count:
            return Response(json.dumps("Not enough Data points"), status=400)

        # Fetches the count of stat_name that exceeds the threshold for given duration
        stats = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(
            DataModel.__getattribute__(DataModel, constants.STAT_REQUEST[stat_name]) > threshold).filter(
            DataModel.date_created > time_obj).filter(DataModel.date_created < time_now).count()

        return jsonify({"ViolatedCount": stats})

    except Exception as e:
        return Response(e, status=404)


# The endpoint returns average of requested stat for a duration, returns error if sufficient data points are not present
@app.route('/stats/avg/<string:stat_name>/<int:duration>')
def average(stat_name, duration):
    # calculate time to query for data 
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points 
    time_obj = time_now - timedelta(minutes=duration)

    stat_list = []
    try:
        # Fetches list of rows that is filter by stat_name and are filtered by decision period
        avg_list = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(
            DataModel.date_created > time_obj).filter(DataModel.date_created < time_now).with_entities(
            text(constants.STAT_REQUEST[stat_name])).all()
        for avg_value in avg_list:
            stat_list.append(avg_value[0])

        # If expected data points count are not present then respond with error
        if duration // sim.frequency_minutes > len(stat_list):
            return Response(json.dumps("Not enough Data points"), status=400)

        # check if any data points were collected
        if not stat_list:
            return Response(json.dumps("Not enough Data points"), status=400)

        # Average, minimum and maximum value of a stat for a given decision period
        return jsonify({
            "avg": sum(stat_list) / len(stat_list),
            "min": min(stat_list),
            "max": max(stat_list), })

    except Exception as e:
        return Response(str(e), status=404)


# The endpoint returns request stat from the latest poll, returns error if sufficient data points are not present.
@app.route('/stats/current/<string:stat_name>')
def current(stat_name):
    try:
        if constants.STAT_REQUEST[stat_name] == constants.CLUSTER_STATE:
            if Simulator.is_provision_in_progress():
                return jsonify({"current_stat": constants.CLUSTER_STATE_YELLOW})
        # Fetches the stat_name for the latest poll
        current_stat = DataModel.query.order_by(desc(DataModel.date_created)).with_entities(
            DataModel.__getattribute__(DataModel, constants.STAT_REQUEST[stat_name])).all()

        # If expected data points count are not present then respond with error
        if len(current_stat) == 0:
            return Response(json.dumps("Not enough Data points"), status=400)

        return jsonify({"current_stat": current_stat[0][constants.STAT_REQUEST[stat_name]]})

    except KeyError:
        return Response(f'stat not found - {stat_name}', status=404)
    except Exception as e:
        return Response(e, status=404)


@app.route('/provision/addnode', methods=["POST"])
def add_node():
    """
    Endpoint to simulate that a cluster state change is under provision
    Expects request body to specify the number of nodes added or removed
    :return: total number of resultant nodes and duration of cluster state as yellow
    """
    try:
        # get the number of added nodes from request body
        nodes = int(request.json['nodes'])
        sim = Simulator(configs.cluster, configs.data_function, configs.searches, configs.simulation_frequency_minutes)
        sim.cluster.add_nodes(nodes)
        cluster_objects = sim.run(24 * 60)

        cluster_objects_post_change = []
        now = datetime.now()
        for cluster_obj in cluster_objects:
            if cluster_obj.date_time >= now:
                cluster_objects_post_change.append(cluster_obj)
                task = DataModel(
                    cpu_usage_percent=cluster_obj.cpu_usage_percent,
                    memory_usage_percent=cluster_obj.memory_usage_percent,
                    date_created=cluster_obj.date_time,
                    total_nodes_count=cluster_obj.total_nodes_count,
                    status=cluster_obj.status
                )
                db.session.merge(task)
        db.session.commit()
        plot_data_points(cluster_objects_post_change, skip_data_ingestion=True)
    except BadRequest as err:
        return Response(json.dumps("expected key 'nodes'"), status=404)
    expiry_time = Simulator.create_provisioning_lock()
    return jsonify({
        'expiry': expiry_time,
        'nodes': sim.cluster.total_nodes_count
    })


@app.route('/all')
def all_data():
    count = DataModel.query.with_entities(
        DataModel.cpu_usage_percent,
        DataModel.memory_usage_percent,
        DataModel.status
    ).count()
    return jsonify(count)


if __name__ == "__main__":
    db.create_all()

    # remove any existing provision lock
    Simulator.remove_provisioning_lock()
    # get configs from config yaml
    configs = parse_config(os.path.join(get_source_code_dir(), constants.CONFIG_FILE_PATH))
    # create the simulator object
    sim = Simulator(configs.cluster, configs.data_function, configs.searches, configs.simulation_frequency_minutes)
    # generate the data points from simulator
    cluster_objects = sim.run(24 * 60)
    # save the generated data points to png
    plot_data_points(cluster_objects)
    # save data points inside db
    for cluster_obj in cluster_objects:
        task = DataModel(
            cpu_usage_percent=cluster_obj.cpu_usage_percent,
            memory_usage_percent=cluster_obj.memory_usage_percent,
            total_nodes_count=cluster_obj.total_nodes_count,
            date_created=cluster_obj.date_time,
            status=cluster_obj.status
        )
        db.session.add(task)
    db.session.commit()

    # start serving the apis
    app.run(port=constants.APP_PORT, debug=True)
