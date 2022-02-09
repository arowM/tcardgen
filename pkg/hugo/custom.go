package hugo

import (
	"io"
	"os"

	"github.com/gohugoio/hugo/parser/pageparser"
)

// in order for codes to pass the original test, i write custom codes instead of overwriting the original codes

func CustomParseFrontMatter(filename string) (*FrontMatter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return customParseFrontMatter(file)
}

func customParseFrontMatter(r io.Reader) (*FrontMatter, error) {
	cfm, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return nil, err
	}

	fm := &FrontMatter{}
	if fm.Title, err = getString(&cfm, fmTitle); err != nil {
		return nil, err
	}
	if isArray := isArray(&cfm, fmAuthor); isArray {
		if fm.Author, err = getFirstStringItem(&cfm, fmAuthor); err != nil {
			return nil, err
		}
	} else {
		if fm.Author, err = wrapGetString(getString(&cfm, fmAuthor)); err != nil {
			return nil, err
		}
	}
	if fm.Category, err = wrapGetFirstStringItem(getFirstStringItem(&cfm, fmCategories)); err != nil {
		return nil, err
	}
	if fm.Tags, err = getAllStringItems(&cfm, fmTags); err != nil {
		return nil, err
	}
	if fm.Date, err = getContentDate(&cfm); err != nil {
		return nil, err
	}

	return fm, nil
}

func wrapGetString(got string, err error) (string, error) {
	if err == nil {
		return got, nil
	}
	if err, ok := err.(*FMNotExistError); !ok {
		return got, err
	} else if err.Key == fmAuthor {
		return "", nil
	}
	return got, err
}

func wrapGetFirstStringItem(got string, err error) (string, error) {
	if err == nil {
		return got, nil
	}
	if err, ok := (err).(*FMNotExistError); !ok {
		return got, err
	} else if err.Key == fmCategories {
		return "", nil
	}
	return got, err
}
