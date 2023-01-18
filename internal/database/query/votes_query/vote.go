package votes_query

import (
	"database/sql"
	"github.com/arandich/telegram-dao/internal/database/entity"
	"log"
)

func AddVote(db *sql.DB, vote *entity.Vote) bool {
	result, err := db.Exec(`INSERT INTO votes (name, url, date_end, text_1, text_2, text_3) VALUES ($1,$2,$3,$4,$5,$6)`, vote.Name, vote.Url, vote.DateEnd, vote.Text1, vote.Text2, vote.Text3)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}

func VoteYes(db *sql.DB, vote *entity.UserVote) bool {
	result, err := db.Exec("UPDATE votes set var_1 = var_1 + $1 WHERE id = $2", vote.Amount, vote.VoteId)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}

func VoteNo(db *sql.DB, vote *entity.UserVote) bool {
	result, err := db.Exec("UPDATE votes set var_2 = var_2 + $1 WHERE id = $2", vote.Amount, vote.VoteId)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}

func VoteElse(db *sql.DB, vote *entity.UserVote) bool {
	result, err := db.Exec("UPDATE votes set var_3 = var_3 + $1 WHERE id = $2", vote.Amount, vote.VoteId)
	if err != nil {
		log.Println(err)
		return false
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return false
	}
	return true
}

func SelectVotes(db *sql.DB, user *entity.User) (*entity.VoteList, bool) {
	rows, err := db.Query(`SELECT * FROM votes WHERE NOT EXISTS(SELECT id FROM votes_journal WHERE votes_journal.user_id = $1 AND votes_journal.vote_id = votes.id)`, user.Id)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer rows.Close()

	votes := entity.VoteList{List: map[string]entity.Vote{}}

	for rows.Next() {
		vote := entity.Vote{}
		err := rows.Scan(&vote.Id, &vote.Name, &vote.Url, &vote.DateStart, &vote.DateEnd, &vote.Text1, &vote.Text2, &vote.Text3, &vote.Var1, &vote.Var2, &vote.Var3)
		if err != nil {
			log.Println(err)
			continue
		}

		votes.List[vote.Name] = vote
	}
	if len(votes.List) == 0 {
		return nil, false
	}

	return &votes, true
}

func SelectVote(db *sql.DB, vote *entity.Vote) (*entity.Vote, bool) {
	rows, err := db.Query(`SELECT * FROM votes WHERE votes.name = $1`, vote.Name)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&vote.Id, &vote.Name, &vote.Url, &vote.DateStart, &vote.DateEnd, &vote.Text1, &vote.Text2, &vote.Text3, &vote.Var1, &vote.Var2, &vote.Var3)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	if vote.Name == "" {
		return nil, false
	}

	return vote, true
}

func SelectVotesArr(db *sql.DB, user *entity.User) (*entity.VoteArr, bool) {
	rows, err := db.Query(`SELECT * FROM votes WHERE NOT EXISTS(SELECT id FROM votes_journal WHERE votes_journal.user_id = $1 AND votes_journal.vote_id = votes.id)`, user.Id)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer rows.Close()

	votes := entity.VoteArr{List: []entity.Vote{}}

	for rows.Next() {
		vote := entity.Vote{}
		err := rows.Scan(&vote.Id, &vote.Name, &vote.Url, &vote.DateStart, &vote.DateEnd, &vote.Text1, &vote.Text2, &vote.Text3, &vote.Var1, &vote.Var2, &vote.Var3)
		if err != nil {
			log.Println(err)
			continue
		}

		votes.List = append(votes.List, vote)
	}
	log.Println(len(votes.List))
	if len(votes.List) == 0 {
		return nil, false
	}

	return &votes, true
}
