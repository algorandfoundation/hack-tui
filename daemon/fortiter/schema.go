package fortiter

var Schema = `
CREATE TABLE IF NOT EXISTS agreements (
    type INTEGER,
    round INTEGER,
    period INTEGER,
    step INTEGER,
  	hash TEXT,
  	sender TEXT,
  	object_round INTEGER,
  	object_period INTEGER,
  	object_step INTEGER,
  	weight INTEGER,
  	weight_total INTEGER,
  	
    message TEXT,
    time INTEGER
);
`
