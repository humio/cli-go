package api

import (
	"strings"

	"github.com/shurcooL/graphql"
)

type ParserTestCase struct {
	Input  string
	Output map[string]string
}

type Parser struct {
	Name      string
	Tests     []ParserTestCase `yaml:",omitempty"`
	Example   string           `yaml:",omitempty"`
	Script    string           `yaml:",flow"`
	TagFields []string         `yaml:",omitempty"`
}

type Parsers struct {
	client *Client
}

func (c *Client) Parsers() *Parsers { return &Parsers{client: c} }

type ParserListItem struct {
	Name      string
	IsBuiltIn bool
}

func (p *Parsers) List(reposistoryName string) ([]ParserListItem, error) {
	var q struct {
		Repository struct {
			Parsers []ParserListItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&q, variables)

	var parsers []ParserListItem
	if graphqlErr == nil {
		parsers = q.Repository.Parsers
	}

	return parsers, graphqlErr
}

func (p *Parsers) Remove(reposistoryName string, parserId string) error {
	var mutation struct {
		RemoveParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeParser(input: { id: $id, repositoryName: $repositoryName })"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
		"id":             graphql.String(parserId),
	}

	return p.client.Mutate(&mutation, variables)
}

func (p *Parsers) Add(reposistoryName string, parser *Parser, force bool) error {

	var mutation struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
	}

	tagFieldsGQL := make([]graphql.String, len(parser.TagFields))

	for i, field := range parser.TagFields {
		tagFieldsGQL[i] = graphql.String(field)
	}

	variables := map[string]interface{}{
		"name":           graphql.String(parser.Name),
		"sourceCode":     graphql.String(parser.Script),
		"repositoryName": graphql.String(reposistoryName),
		"testData":       testCasesToStrings(parser),
		"tagFields":      tagFieldsGQL,
		"force":          graphql.Boolean(force),
	}

	return p.client.Mutate(&mutation, variables)
}

func testCasesToStrings(parser *Parser) []graphql.String {

	lines := strings.Split(parser.Example, "\n")

	result := make([]graphql.String, 0)
	for _, item := range parser.Tests {
		result = append(result, graphql.String(item.Input))
	}

	for i, item := range lines {
		if i != len(lines)-1 {
			result = append(result, graphql.String(item))
		}
	}

	return result
}

func (p *Parsers) Get(reposistoryName string, parserId string) (*Parser, error) {

	var query struct {
		Repository struct {
			Parser struct {
				Name       string
				SourceCode string
				TestData   []string
				TagFields  []string
			} `graphql:"parser(id: $parserId)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserId":       graphql.String(parserId),
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&query, variables)

	var parser Parser
	if graphqlErr == nil {
		parser = Parser{
			Name:      query.Repository.Parser.Name,
			Tests:     mapTests(query.Repository.Parser.TestData, toTestCase),
			Script:    query.Repository.Parser.SourceCode,
			TagFields: query.Repository.Parser.TagFields,
		}
	}

	return &parser, graphqlErr
}

func mapTests(vs []string, f func(string) ParserTestCase) []ParserTestCase {
	vsm := make([]ParserTestCase, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func toTestCase(line string) ParserTestCase {
	return ParserTestCase{
		Input:  line,
		Output: map[string]string{},
	}
}

func (p *Parsers) Export(reposistoryName string, parserId string) (string, error) {

	var query struct {
		Repository struct {
			Parser struct {
				YamlTemplate string
			} `graphql:"parser(id: $parserId)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserId":       graphql.String(parserId),
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&query, variables)

	if graphqlErr != nil {
		return "", graphqlErr
	}

	return query.Repository.Parser.YamlTemplate, nil
}
