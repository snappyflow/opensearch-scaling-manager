from flask import Flask,jsonify,Response
from flask_sqlalchemy import SQLAlchemy
from datetime import datetime,timedelta
from sqlalchemy import func,text,desc
from simulator import Simulator
from cluster import Cluster
from data_ingestion import State, DataIngestion
import constants
import json
import os
from config_parser import parse_config
import shutil

DATA_POINT_FREQUENCY_MINUTES = 5 
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
    date_created = db.Column(db.DateTime, default = datetime.now(),primary_key =True)


# Converts the duration in minutes to time object of "HH:MM" format
def convert_to_hh_mm(duration_in_m):
    time_h_m  = '{:02d}:{:02d}'.format(*divmod(duration_in_m, 60))
    time_obj  = datetime.strptime(time_h_m, '%H:%M')
    return time_obj



# Returns the violated count for a requested metric, threshold and duration, returns error if sufficient data points are not present.
@app.route('/stats/violated/<string:stat_name>/<int:duration>/<float:threshold>')
def violated_count(stat_name, duration, threshold):
    # calculate time to query for data 
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points
    time_obj = time_now - timedelta(minutes=duration)

    try:
        # Fetching the count of data points for given duration.
        data_point_count = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(DataModel.date_created > time_obj).count()

        # If expected data points are not present then respond with error
        if duration/DATA_POINT_FREQUENCY_MINUTES > data_point_count:
            return Response("Not enough data points",status=400)

        # Fetches the count of stat_name that exceeds the threshold for given duration
        stats = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(DataModel.cpu_usage_percent > threshold).filter(DataModel.date_created > time_obj).count()
        
        return jsonify({"ViolatedCount" : stats})
    
    except Exception as e:
        return Response(e,status=404)

# The endpoint returns average of requested stat for a duration, returns error if sufficient data points are not present
@app.route('/stats/avg/<string:stat_name>/<int:duration>')
def average(stat_name,duration):
    # calculate time to query for data 
    time_now = datetime.now()

    # Convert the minutes to time object to compare and query for required data points 
    time_obj = time_now - timedelta(minutes=duration)

    stat_list = []
    try:
        # Fetches list of rows that is filter by stat_name and are filterd by decision period
        avg_list = DataModel.query.order_by(constants.STAT_REQUEST[stat_name]).filter(DataModel.date_created > time_obj).with_entities(text(constants.STAT_REQUEST[stat_name])).all()
        for i in avg_list:  
            stat_list.append(i[0])

        # If expected data points count are not present then respond with error
        if duration/DATA_POINT_FREQUENCY_MINUTES > len(stat_list):
            return Response("Not enough data points",status=400)

        # Average, minimum and maximum value of a stat for a given decision period
        return jsonify({
            "avg": sum(stat_list)/len(stat_list),
            "min": min(stat_list),
            "max": max(stat_list),})

    except Exception as e:
        return Response(e,status=404)


# The endpoint returns request stat from the latest poll, returns error if sufficient data points are not present.
@app.route('/stats/current/<string:stat_name>')
def current(stat_name):
    try:
        # Fetches the stat_name for the lastest poll  
        current = DataModel.query.order_by(desc(DataModel.date_created)).with_entities(DataModel.__getattribute__(DataModel,constants.STAT_REQUEST[stat_name])).all()
        
        # If expected data points count are not present then respond with error
        if len(current) == 0:
            return Response("Not enough Data points",status=400)

        return jsonify({"current": current[0][constants.STAT_REQUEST[stat_name]]})

    except Exception as e:
        return Response(e,status=404)

@app.route('/all')
def all():
    task = DataModel.query.with_entities(DataModel.cpu_usage_percent,DataModel.memory_usage_percent,DataModel.status,DataModel.date_created).count()
    return jsonify(task)
    

if __name__ == "__main__":
    db.create_all()
    configs = parse_config('config.yaml')
    all_states = [State(**state) for state in configs.data_ingestion.get('states')]
    randomness_percentage = configs.data_ingestion.get('randomness_percentage')

    data_function = DataIngestion(all_states, randomness_percentage)

    cluster = Cluster(**configs.stats)

    sim = Simulator(cluster, data_function, configs.searches, configs.data_generation_interval_minutes, 0)
    # generate the data points from simulator
    output = sim.run(24*60)
    # start serving the apis
    for sim_obj,timestamp in output:
        task = DataModel(cpu_usage_percent = sim_obj.cpu_usage_percent,
            memory_usage_percent = sim_obj.memory_usage_percent,
            date_created = timestamp,
            status=sim_obj.status
        )
        db.session.add(task)
    db.session.commit()       
    app.run(port=constants.APP_PORT,debug=True,use_reloader=False)

