from prometheus_client import generate_latest, Counter
from hashlib import sha256
from time import sleep
from random import randint
from flask import Flask, Response, request
from logging.config import dictConfig

app = Flask(__name__)

ROLL_COUNTER = Counter(
        'dice_roll_count', 'How often the dice was rolled'
)

NUMBERS_COUNTER = Counter(
        'dice_numbers_count', 'How often each number of the dice was rolled',
        ['number']
)

@app.route("/rolldice")
def roll_dice():
    player = request.args.get('player', default="Anonymous player")
    max = 8 if sha256(bytes(player, 'utf-8')).hexdigest() == 'f4b7c19317c929d2a34297d6229defe5262fa556ef654b600fc98f02c6d87fdc' else 6
    result = str(do_roll(max))
    app.logger.info("%s is rolling the dice: %s", player, result);
    ROLL_COUNTER.inc()
    NUMBERS_COUNTER.labels(result).inc()
    return result

def do_roll(max):
    result = randint(1, max)
    if result > 6:
        sleep(0.1 * result)
    return result

@app.route("/metrics/")
def metrics():
    return Response(generate_latest(), mimetype=str('text/plain; version=0.0.4; charset=utf-8'))
