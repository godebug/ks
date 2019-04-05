package history

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type History struct {
	Items []*item
}

type item struct {
	Time   time.Time
	Series int64
	Id     int64
	Values []int64
	Value  int64
	Stats  map[key]*value
}

type key struct {
	Max  int64
	Side int64
}

type value struct {
	Qnty   int64
	Volume float64
}

func NewHistory(file string) (*History, error) {
	out := History{}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	cnt := 1
	var past int64
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		i, err := makeItem(record)
		if err != nil {
			return nil, err
		}
		cnt++
		if past != i.Series {
			cnt = 1
		}
		past = i.Series
		i.Time = time.Date(i.Time.Year(), i.Time.Month(), i.Time.Day(), 0, 0, 0, 0, time.UTC)
		i.Time = i.Time.Add(time.Duration(cnt*5) * time.Minute)
		log.Println(i.Series, i.Id, i.Time.Format("02.01.2006 15:04:05"))

		out.Items = append(out.Items, i)
	}
	return &out, nil
}

func makeItem(in []string) (out *item, err error) {
	out = &item{}
	out.Time, err = time.Parse("02.01.2006", in[0])
	if err != nil {
		return nil, err
	}

	i64, err := strconv.ParseInt(in[1], 10, 64)
	if err != nil {
		return nil, err
	}
	out.Id = i64
	i64, err = strconv.ParseInt(in[2], 10, 64)
	if err != nil {
		return nil, err
	}
	out.Series = i64

	for i := 3; i < 11; i++ {
		i64, err = strconv.ParseInt(in[i], 10, 64)
		if err != nil {
			return nil, err
		}
		out.Values = append(out.Values, i64)
	}
	i64, err = strconv.ParseInt(in[11], 10, 64)
	if err != nil {
		return nil, err
	}
	out.Value = i64
	out.Stats = make(map[key]*value)
	const pos = 12
	for i := 0; i < 10; i++ {
		idx := i * 3
		k, err := makeKey(in[pos+idx])
		if err != nil {
			return nil, err
		}
		v, err := makeValue(in[pos+idx+1], in[pos+idx+2])
		if err != nil {
			return nil, err
		}
		out.Stats[k] = v
	}
	return out, nil
}

func makeKey(raw string) (out key, err error) {
	ab := strings.Split(raw, "/")
	if len(ab) != 2 {
		return out, errors.New("Wrong key: " + raw)
	}
	i64, err := strconv.ParseInt(ab[0], 10, 64)
	if err != nil {
		return out, err
	}
	out.Max = i64

	i64, err = strconv.ParseInt(ab[1], 10, 64)
	if err != nil {
		return out, err
	}
	out.Side = i64
	return out, nil
}

func makeValue(a, b string) (out *value, err error) {
	i64, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return out, err
	}
	f64, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return out, err
	}
	return &value{i64, f64}, nil
}
