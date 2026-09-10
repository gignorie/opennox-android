package cfg

import (
	"bufio"
	"io"
	"strings"
)

type KeyValue struct {
	Key     string
	Value   string
	Comment string
}

type Section []KeyValue

func (s Section) Get(key string) (string, bool) {
	for _, kv := range s {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return "", false
}

func (s *Section) Set(key, val string) {
	for i, kv := range *s {
		if kv.Key == key {
			(*s)[i].Value = val
			return
		}
	}
	*s = append(*s, KeyValue{Key: key, Value: val})
}

type File struct {
	Sections []Section
}

func Parse(r io.Reader) (*File, error) {
	sc := bufio.NewScanner(r)
	f := &File{}
	var cur Section
	add := func(kv KeyValue) {
		if kv == (KeyValue{}) {
			return
		}
		cur = append(cur, kv)
	}
	done := func(force bool) {
		if force || len(cur) != 0 {
			f.Sections = append(f.Sections, cur)
			cur = nil
		}
	}
	defer done(false)
	comment := ""
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(line[1:])
			if comment == "" {
				comment = line
			} else {
				comment += "\n" + line
			}
			continue
		}
		kv := KeyValue{Comment: comment}
		if strings.HasPrefix(line, "---") {
			add(kv)
			done(true)
			continue
		}
		comment = ""
		sub := strings.SplitN(line, "=", 2)
		kv.Key = strings.TrimSpace(sub[0])
		if len(sub) == 2 {
			kv.Value = strings.TrimSpace(sub[1])
		}
		add(kv)
	}
	return f, sc.Err()
}

func (f *File) WriteTo(w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, sect := range f.Sections {
		for _, kv := range sect {
			if kv.Comment != "" {
				for _, c := range strings.Split(kv.Comment, "\n") {
					bw.WriteString("# ")
					bw.WriteString(c)
					bw.WriteString("\n")
				}
			}
			bw.WriteString(strings.TrimSpace(kv.Key))
			bw.WriteString(" = ")
			bw.WriteString(strings.TrimSpace(kv.Value))
			bw.WriteString("\n")
		}
		bw.WriteString("---\n")
	}
	return bw.Flush()
}
