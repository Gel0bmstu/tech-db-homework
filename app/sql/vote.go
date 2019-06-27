package slc

var (
	CreateVote = `INSERT INTO votes
				  (thread, nickname, voice)
				  VALUES ($1, $2, $3)
				  ON CONFLICT (nickname, thread)
				  DO UPDATE SET voice = excluded.voice;`
)
