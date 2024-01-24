package main

func main() {
	err := AddFile()
	if err != nil {
		// log.Error(err)
	}
}

func AddFile() error {
	// c, err := NewClient()
	// if err != nil {
	// 	return err
	// }
	// // add and pin file
	// ctx := context.Background()
	// sth := test.NewShardingTestHelper()
	// mfr, closer := sth.GetTreeMultiReader(&testing.T{})
	// defer closer.Close()
	// pin, err := c.AddECMultiFile(ctx, mfr, "shardTesting")
	// if err != nil {
	// 	return err
	// }

	// // check if file is pinned
	// c.GetFileDag(pin.Cid.String())
	return nil
}
