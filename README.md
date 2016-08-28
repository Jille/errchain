# errchain

[![Build Status](https://travis-ci.org/Jille/errchain.png)](https://travis-ci.org/Jille/errchain)

errchain is a simple library helping you with error handling.

https://godoc.org/github.com/Jille/errchain

Usage example:
```go
func write(data T) (retErr error) {
	fh, err := os.Create("data.gz.new")
	if err != nil {
		return err
	}
	gz := gzip.NewWriter(fh)
	enc := gob.NewEncoder(gz)
	err = enc.Encode(data)
	errchain.Append(&err, gz.Close())
	errchain.Append(&err, fh.Close())
	if err != nil {
		return err
	}
	return os.Rename("data.gz.new", "data.gz")
}

func read() (data T, retErr error) {
	fh, err := os.Open("data.gz")
	if err != nil {
		return nil, err
	}
	defer errchain.Call(&retErr, fh.Close)
	gz, err := gzip.NewReader(fh)
	if err != nil {
		return nil, err
	}
	defer errchain.Call(&retErr, gz.Close)
	dec := gob.NewDecoder(gz)
	if err := dec.Decode(data); err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	data, err := read()
	if err == nil {
		data++
		err = write(data)
	}
	if err != nil {
		errs := errchain.List(err)
		if len(errs) > 1 {
			fmt.Printf("%d errors:\n", len(errs))
		}
		for _, e := range errs {
			fmt.Println(e)
		}
	}
}
```
