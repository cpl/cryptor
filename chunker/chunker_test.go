package chunker_test

// func TestChunker(t *testing.T) {
// 	t.Parallel()

// 	db, err := ldbcache.New("/tmp/cryptordb", 0, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.RemoveAll("/tmp/cryptordb")
// 	manager := cachedb.New("/tmp/cryptordb", db)

// 	c := chunker.New(con.MB, manager)
// 	archive.TarGz("../", c)

// 	tail, err := c.Pack(aes.NewKey())
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	log.Println(manager.Count(), manager.Size())
// 	i := manager.Iterator()
// 	for i.Next() {
// 		log.Println(b16.EncodeString(i.Key()))
// 	}

// 	log.Println(b16.EncodeString(tail))
// }
