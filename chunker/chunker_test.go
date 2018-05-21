package chunker_test

// func TestChunker(t *testing.T) {
// 	t.Parallel()

// 	db, err := ldbcache.New("/tmp/cryptor_db", 0, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	c := chunker.New(con.KB, db)
// 	archive.TarGz(".", c)

// 	if err := c.Pack(aes.NewKey()); err != nil {
// 		t.Error(err)
// 	}
// }
