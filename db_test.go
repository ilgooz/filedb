package filedb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDial(t *testing.T) {

	db, err := Dial("test/missing")
	require.Equal(t, err, ErrDBNotFound)
	require.Nil(t, db)

	db, err = Dial("test/file.txt")
	require.Equal(t, err, ErrDBNotFound)
	require.Nil(t, db)

	db, err = Dial("test/db")
	require.NoError(t, err)
	require.NotNil(t, db)
	db.Close()

}

func TestC(t *testing.T) {

	db, err := Dial("test/db")
	require.NoError(t, err)
	defer db.Close()
	c, err := db.C("people")
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, c.Path(), "test/db/people.filedb")
	require.Equal(t, c.DB(), db)

	// same C
	c2, _ := db.C("people")
	require.Equal(t, c, c2, "Cs with same name should be same object")

}

func TestCollections(t *testing.T) {

	db, err := Dial("test/db")
	require.NoError(t, err)
	defer db.Close()

	c2, _ := db.C("events")
	c3, _ := db.C("comments")

	c2.Insert([]byte("something"))
	c3.Insert([]byte("something"))

	cols, err := db.ColNames()
	require.NoError(t, err)

	require.Equal(t, 3, len(cols))
	require.Equal(t, cols[0], "comments")
	require.Equal(t, cols[1], "events")
	require.Equal(t, cols[2], "people")

}

func TestInsert(t *testing.T) {

	db, err := Dial("test/db")
	require.NoError(t, err)
	defer db.Close()
	c, err := db.C("people")

	obj := map[string]interface{}{"name": "Mat", "location": "San Francisco"}
	b, _ := json.Marshal(&obj)
	err = c.Insert(b)
	require.NoError(t, err)

}

func TestForEach(t *testing.T) {

	db, err := Dial("test/db")
	require.NoError(t, err)
	defer db.Close()
	c, err := db.C("people")

	require.NoError(t, c.Drop())

	require.NoError(t, c.InsertJSON(map[string]interface{}{"name": "Mat", "location": "San Francisco"}))
	require.NoError(t, c.InsertJSON(map[string]interface{}{"name": "Ryan", "location": "Boulder"}))
	require.NoError(t, c.InsertJSON(map[string]interface{}{"name": "Tyler", "location": "Salt Lake City"}))

	var lines []string
	err = c.ForEach(func(j []byte) bool {
		lines = append(lines, string(j))
		return false // don't break
	})
	require.NoError(t, err)
	require.Equal(t, len(lines), 3)

}
