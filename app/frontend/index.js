/*
 Put your NodeSDK initialization here.
 */
const { context, trace, metrics, ValueType } = require('@opentelemetry/api');
const http = require("http");
const app = require("express")();
const pino = require('pino-http')()

// not included by default
//var otelsdkinit = require('./instrument.js');

app.use(pino)

const port = process.env.FRONTEND_PORT || 4000;
const backend1url =
  process.env.BACKEND1_URL || "http://localhost:5165/rolldice";
const backend2url =
  process.env.BACKEND2_URL || "http://localhost:5000/rolldice";

const myMeter = metrics.getMeter("app-meter");
const requestCounter = myMeter.createCounter('request_total', {
  description: "Counter of requests",
  valueType: ValueType.INT
})
const gameCounter = myMeter.createUpDownCounter('app_games_total', {
  description: "A counter of how often the game has been played",
  valueType: ValueType.INT
})
const winCounter = myMeter.createUpDownCounter('app_wins_total', {
  description: "A counter per player who has won",
  valueType: ValueType.INT
})

app.get("/", (req, res) => {
  requestCounter.add(1);
  const { player1, player2 } = Object.assign({player1: "Player 1", player2: "Player 2"}, req.query)
  if(player1 == 'Player 1') {
    req.log.info('Player 1 prefers to stay anonymous.')
  }
  if(player2 == 'Player 2') {
    req.log.info('Player 2 prefers to stay anonymous.')
  }
  span = trace.getSpan(context.active())
  if(span) {
    span.setAttribute('app.player1', player1)
    // TODO(tracing): Add an attribute for player2
  }

  const p1 = new Promise((resolve, reject) => {
    http.get(`${backend1url}?player=${player1}`, (response) => {
      let data = [];

      response.on("data", (chunk) => {
        data.push(chunk);
      });
      response.on("end", () => {
        try {
          const result = JSON.parse(Buffer.concat(data).toString());
          res.write("Player 1 rolls: " + result + "\n");
          resolve(result);
        } catch(error) {
          reject(error)
        }
      });
    }).on('error', (error) => {
        req.log.error("Backend1 is not available.")
        reject(error)
    }).end()
  });

  const p2 = new Promise((resolve, reject) => {
    http.get(`${backend2url}?player=${player2}`, (response) => {
      let data = [];

      response.on("data", (chunk) => {
        data.push(chunk);
      });
      response.on("end", () => {
        try {
          const result = Buffer.concat(data).toString();
          res.write("Player 2 rolls: " + result + "\n");
          resolve(result);
        } catch(error) {
          reject(error)
        }
      });
    }).on('error', (error) => {
      req.log.error("Backend2 is not available.")
      reject(error)
    }).end()
  });

  Promise.all([p1, p2]).then(([roll1, roll2]) => {
    let winner = 'Nobody'
    let winnerRolled = 0
    if (roll1 > roll2) {
      winner = player1
      winnerRolled = roll1
    } else if (roll2 > roll1) {
      winner = player2
      winnerRolled = roll2
    }
    // TODO(tracing): Add the winner as a span attribute

    // Count the total number of games
    gameCounter.add(1);

    // TODO (metrics): count how often each player wins - winCounter
    // use app.winner tag
    // Count how often each player wins

    // Add counters for numbers rolled and/or for players who played

    res.end(`${winner} wins`);
  }).catch(error => {
    try {
      res.sendStatus(500).end()
    } catch(e) {
      // If sending the error fails, the service crashes, we want to avoid that!
    }
  });
});

app.listen(port, () => {
  pino.logger.info(`Example app listening on port ${port}`);
});
