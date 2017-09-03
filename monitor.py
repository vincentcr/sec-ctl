#!/usr/bin/env python

import logging

from flask import Flask, request, jsonify
from jsonschema import validate, ValidationError
from envisalink import Envisalink, ClientCodes
from MessageJSONEncoder import MessageJSONEncoder
from Config import Config


app = Flask("alarm-monitor")
app.json_encoder = MessageJSONEncoder
app.use_reloader=False

@app.route("/")
def root():
    return "hello, worldas"


@app.route("/cmd", methods=['POST'])
def parse_cmd():
    schema = {
        'type': 'object',
        'properties': {
            "code": {"type": "string"},
            "data": {"type": "string"},
        },
        "required": ["code"],
        "additionalProperties": False,
    }
    json = request.get_json()
    try:
        validate(json, schema)
    except ValidationError as e:
        print("bad input", e)
        return "Bad input", 400

    code = ClientCodes.by_name(json['code'])
    data = json.get('data')
    (cmd_id, reply) = tpi.send_cmd(code, data=data)
    print('code=', code, 'data=', data, 'cmd_id', cmd_id, 'reply=', reply)
    return jsonify(reply)


if __name__ == "__main__":
    global tpi, config
    logging.getLogger().setLevel(logging.DEBUG)
    config = Config.load()
    tpi = Envisalink(
        hostname=config['envisalink']['host'],
        password=config['envisalink']['password'],
        )
    app.run(
        host='0.0.0.0',
        port=config['http']['port']
    )
