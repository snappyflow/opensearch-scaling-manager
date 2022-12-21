import os
import shutil
import json
from datetime import datetime, timedelta

from flask import Flask, jsonify, Response, request
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import text, desc

import constants
from config_parser import parse_config, get_source_code_dir
from data_ingestion import State, DataIngestion
from cluster import Cluster
from open_search_simulator import Simulator
from plotter import plot_data_points


app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///datapoints.db"
app.app_context().push()
if os.path.exists("instance"):
    shutil.rmtree("instance")
db = SQLAlchemy(app)


# Database model to store the datapoints
class DataModel(db.Model):
    status = db.Column(db.String(200))
    cpu_usage_percent = db.Column(db.Float, default=0)
    memory_usage_percent = db.Column(db.Float, default=0)
    shards_count = db.Column(db.Integer, default=0)
    date_created = db.Column(db.DateTime, default=datetime.now(), primary_key=True)


def convert_to_hh_mm(duration_in_m):
    """
    Converts the duration in minutes to time object of "HH:MM" format
    :param duration_in_m: represnts the duration in minutes 
    :return: Time object generated from minutes
    """
    time_h_m = "{:02d}:{:02d}".format(*divmod(duration_in_m, 60))
    time_obj = datetime.strptime(time_h_m, "%H:%M")
    return time_obj


@app.route("/stats/violated/<string:stat_name>/<int:duration>/<float:threshold>")
def violated_count(stat_name, duration, threshold):
    """
    Endpoint fetches the violated count for a requested metric, threshold and duration,
    :param stat_name: represents the stat that is being queried.
    :param duration: represents the time period for fetaching the average
    :param threshold: represents the limit considered for evaluating violated count
    :return: count of stat exceeding the threshold for a given duration
    """
    # calculate time to query for data
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points
    query_begin_time = time_now - timedelta(minutes=duration)

    try:
        # Fetching the count of data points for given duration.
        data_point_count = (
            DataModel.query.order_by(constants.STAT_REQUEST[stat_name])
            .filter(DataModel.date_created > query_begin_time)
            .filter(DataModel.date_created < time_now)
            .count()
        )

        # If expected data points are not present then respond with error
        if first_data_point_time > query_begin_time:
            return Response(json.dumps("Not enough Data points"), status=400)

        # Fetches the count of stat_name that exceeds the threshold for given duration
        stats = (
            DataModel.query.order_by(constants.STAT_REQUEST[stat_name])
            .filter(
                DataModel.__getattribute__(DataModel, constants.STAT_REQUEST[stat_name])
                > threshold
            )
            .filter(DataModel.date_created > query_begin_time)
            .filter(DataModel.date_created < time_now)
            .count()
        )

        return jsonify({"ViolatedCount": stats})

    except Exception as e:
        return Response(e, status=404)


@app.route("/stats/avg/<string:stat_name>/<int:duration>")
def average(stat_name, duration):
    """
    The endpoint cevaluates average of requested stat for a duration
    returns error if sufficient data points are not present.
    :param stat_name: represents the stat that is being queried.
    :param duration: represents the time period for fetaching the average
    :return: average of the provided stat name for the decision period.
    """
    # calculate time to query for data
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points
    query_begin_time = time_now - timedelta(minutes=duration)

    stat_list = []
    try:
        # Fetches list of rows that is filter by stat_name and are filterd by decision period
        avg_list = (
            DataModel.query.order_by(constants.STAT_REQUEST[stat_name])
            .filter(DataModel.date_created > query_begin_time)
            .filter(DataModel.date_created < time_now)
            .with_entities(text(constants.STAT_REQUEST[stat_name]))
            .all()
        )
        for avg_value in avg_list:
            stat_list.append(avg_value[0])

        # If expected data points count are not present then respond with error
        if first_data_point_time > query_begin_time:
            return Response(json.dumps("Not enough Data points"), status=400)

        # check if any data points were collected
        if not stat_list:
            return Response(json.dumps("Not enough Data points"), status=400)

        # Average, minimum and maximum value of a stat for a given decision period
        return jsonify(
            {
                "avg": sum(stat_list) / len(stat_list),
                "min": min(stat_list),
                "max": max(stat_list),
            }
        )

    except Exception as e:
        return Response(str(e), status=404)


@app.route("/stats/current/<string:stat_name>")
def current(stat_name):
    """
    The endpoint to fetch stat from the latest poll,
    Returns error if sufficient data points are not present.
    :return: Stat generated by the most recent poll
    """
    try:
        if constants.STAT_REQUEST[stat_name] == constants.CLUSTER_STATE:
            if Simulator.is_provision_in_progress():
                return jsonify({"current": constants.CLUSTER_STATE_YELLOW})
        # Fetches the stat_name for the latest poll
        current = (
            DataModel.query.order_by(desc(DataModel.date_created))
            .with_entities(
                DataModel.__getattribute__(DataModel, constants.STAT_REQUEST[stat_name])
            )
            .all()
        )

        # If expected data points count are not present then respond with error
        if len(current) == 0:
            return Response(json.dumps("Not enough Data points"), status=400)

        return jsonify({"current": current[0][constants.STAT_REQUEST[stat_name]]})

    except Exception as e:
        return Response(str(e), status=404)


@app.route("/provision/addnode", methods=["POST"])
def add_node():
    """
    Endpoint to simulate that a cluster state change is under provision
    Expects request body to specify the number of nodes added or removed
    :return: total number of resultant nodes and duration of cluster state as yellow
    """
    try:
        request.json["nodes"]
    except:
        return Response(json.dumps("Not enough Data points"), status=404)
    # Todo - Reflect node count in cluster
    expiry_time = Simulator.create_provisioning_lock()
    return jsonify({"expiry": expiry_time})


@app.route("/all")
def all():
    """
    Endpoint to fetch the count of datapoints generated by simulator
    :return: returns the count of total datapoints generated by the simulator
    """
    task = DataModel.query.with_entities(
        DataModel.cpu_usage_percent, DataModel.memory_usage_percent, DataModel.status
    ).count()
    return jsonify(task)


if __name__ == "__main__":
    db.create_all()

    # remove any existing provision lock
    Simulator.remove_provisioning_lock()

    configs = parse_config(
        os.path.join(get_source_code_dir(), constants.CONFIG_FILE_PATH)
    )
    all_states = [
        State(**state)
        for state in configs.data_ingestion.get(constants.DATA_INGESTION_STATES)
    ]
    randomness_percentage = configs.data_ingestion.get(
        constants.DATA_INGESTION_RANDOMNESS_PERCENTAGE
    )

    data_function = DataIngestion(all_states, randomness_percentage)

    cluster = Cluster(**configs.stats)

    sim = Simulator(
        cluster,
        data_function,
        configs.searches,
        configs.data_generation_interval_minutes,
    )
    # generate the data points from simulator
    cluster_objects = sim.run(24 * 60)
    first_data_point_time = cluster_objects[0].date_time
    plot_data_points(cluster_objects)
    for cluster_obj in cluster_objects:
        task = DataModel(
            cpu_usage_percent=cluster_obj.cpu_usage_percent,
            memory_usage_percent=cluster_obj.memory_usage_percent,
            date_created=cluster_obj.date_time,
            status=cluster_obj.status,
        )
        db.session.add(task)
    db.session.commit()
    # start serving the apis
    app.run(port=constants.APP_PORT, debug=True)
