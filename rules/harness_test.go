package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/graphql-go/graphql/testutil"
	"reflect"
)

var beingInterface = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Being",
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
	},
})
var petInterface = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Pet",
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
	},
})
var dogCommandEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "DogCommand",
	Values: graphql.EnumValueConfigMap{
		"SIT": &graphql.EnumValueConfig{
			Value: 0,
		},
		"HEEL": &graphql.EnumValueConfig{
			Value: 1,
		},
		"DOWN": &graphql.EnumValueConfig{
			Value: 2,
		},
	},
})
var dogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dog",
	IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
		return true
	},
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
		"nickname": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"barkVolume": &graphql.FieldConfig{
			Type: graphql.Int,
		},
		"barks": &graphql.FieldConfig{
			Type: graphql.Boolean,
		},
		"doesKnowCommand": &graphql.FieldConfig{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"dogCommand": &graphql.ArgumentConfig{
					Type: dogCommandEnum,
				},
			},
		},
		"isHousetrained": &graphql.FieldConfig{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"atOtherHomes": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: true,
				},
			},
		},
		"isAtLocation": &graphql.FieldConfig{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"x": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"y": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
		},
	},
	Interfaces: []*graphql.Interface{
		beingInterface,
		petInterface,
	},
})
var furColorEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "FurColor",
	Values: graphql.EnumValueConfigMap{
		"BROWN": &graphql.EnumValueConfig{
			Value: 0,
		},
		"BLACK": &graphql.EnumValueConfig{
			Value: 1,
		},
		"TAN": &graphql.EnumValueConfig{
			Value: 2,
		},
		"SPOTTED": &graphql.EnumValueConfig{
			Value: 3,
		},
	},
})

var catType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Cat",
	IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
		return true
	},
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
		"nickname": &graphql.FieldConfig{
			Type: graphql.String,
		},
		"meowVolume": &graphql.FieldConfig{
			Type: graphql.Int,
		},
		"meows": &graphql.FieldConfig{
			Type: graphql.Boolean,
		},
		"furColor": &graphql.FieldConfig{
			Type: furColorEnum,
		},
	},
	Interfaces: []*graphql.Interface{
		beingInterface,
		petInterface,
	},
})
var catOrDogUnion = graphql.NewUnion(graphql.UnionConfig{
	Name: "CatOrDog",
	Types: []*graphql.Object{
		dogType,
		catType,
	},
	ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
		// not used for validation
		return nil
	},
})
var intelligentInterface = graphql.NewInterface(graphql.InterfaceConfig{
	Name: "Intelligent",
	Fields: graphql.FieldConfigMap{
		"iq": &graphql.FieldConfig{
			Type: graphql.Int,
		},
	},
})

var humanType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Human",
	IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
		return true
	},
	Interfaces: []*graphql.Interface{
		beingInterface,
		intelligentInterface,
	},
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
		"pets": &graphql.FieldConfig{
			Type: graphql.NewList(petInterface),
		},
		"iq": &graphql.FieldConfig{
			Type: graphql.Int,
		},
		// `relatives` field added later in init()
	},
})

var alienType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Alien",
	IsTypeOf: func(value interface{}, info graphql.ResolveInfo) bool {
		return true
	},
	Interfaces: []*graphql.Interface{
		beingInterface,
		intelligentInterface,
	},
	Fields: graphql.FieldConfigMap{
		"name": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"surname": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
		"iq": &graphql.FieldConfig{
			Type: graphql.Int,
		},
		"numEyes": &graphql.FieldConfig{
			Type: graphql.Int,
		},
	},
})
var dogOrHumanUnion = graphql.NewUnion(graphql.UnionConfig{
	Name: "DogOrHuman",
	Types: []*graphql.Object{
		dogType,
		humanType,
	},
	ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
		// not used for validation
		return nil
	},
})
var humanOrAlienUnion = graphql.NewUnion(graphql.UnionConfig{
	Name: "HumanOrAlien",
	Types: []*graphql.Object{
		alienType,
		humanType,
	},
	ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
		// not used for validation
		return nil
	},
})

var complexInputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ComplexInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"requiredField": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
		"intField": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"stringField": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"booleanField": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"stringListField": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
		},
	},
})
var complicatedArgs = graphql.NewObject(graphql.ObjectConfig{
	Name: "ComplicatedArgs",
	// TODO List
	// TODO Coercion
	// TODO NotNulls
	Fields: graphql.FieldConfigMap{
		"intArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"intArg": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
		},
		"nonNullIntArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"nonNullIntArg": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		},
		"stringArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"stringArg": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
		},
		"booleanArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"booleanArg": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
		},
		"enumArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"enumArg": &graphql.ArgumentConfig{
					Type: furColorEnum,
				},
			},
		},
		"floatArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"floatArg": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
		},
		"idArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"idArg": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
			},
		},
		"stringListArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"stringListArg": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
		"complexArgField": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"complexArg": &graphql.ArgumentConfig{
					Type: complexInputObject,
				},
			},
		},
		"multipleReqs": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"req1": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"req2": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		},
		"multipleOpts": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"opt1": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"opt2": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
		},
		"multipleOptAndReq": &graphql.FieldConfig{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"req1": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"req2": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"opt1": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
				"opt2": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
		},
	},
})
var queryRoot *graphql.Object
var defaultSchema *graphql.Schema

func init() {

	humanType.AddFieldConfig("relatives", &graphql.FieldConfig{
		Type: graphql.NewList(humanType),
	})
	queryRoot = graphql.NewObject(graphql.ObjectConfig{
		Name: "QueryRoot",
		Fields: graphql.FieldConfigMap{
			"human": &graphql.FieldConfig{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Type: humanType,
			},
			"alien": &graphql.FieldConfig{
				Type: alienType,
			},
			"dog": &graphql.FieldConfig{
				Type: dogType,
			},
			"cat": &graphql.FieldConfig{
				Type: catType,
			},
			"pet": &graphql.FieldConfig{
				Type: petInterface,
			},
			"catOrDog": &graphql.FieldConfig{
				Type: catOrDogUnion,
			},
			"dogOrHuman": &graphql.FieldConfig{
				Type: dogOrHumanUnion,
			},
			"humanOrAlien": &graphql.FieldConfig{
				Type: humanOrAlienUnion,
			},
			"complicatedArgs": &graphql.FieldConfig{
				Type: complicatedArgs,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryRoot,
	})
	if err != nil {
		panic(err)
	}
	defaultSchema = &schema

}
func expectValid(t *testing.T, schema *graphql.Schema, rules []graphql.ValidationRuleFn, queryString string) {
	source := source.NewSource(&source.Source{
		Body: queryString,
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.ValidateDocument(schema, AST, rules)
	if len(result.Errors) > 0 {
		t.Fatalf("Should validate, got %v", result.Errors)
	}
	if result.IsValid != true {
		t.Fatalf("IsValid should be true, got %v", result.IsValid)
	}

}
func expectInvalid(t *testing.T, schema *graphql.Schema, rules []graphql.ValidationRuleFn, queryString string, expectedErrors []gqlerrors.FormattedError) {
	source := source.NewSource(&source.Source{
		Body: queryString,
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.ValidateDocument(schema, AST, rules)
	if len(result.Errors) != len(expectedErrors) {
		t.Fatalf("Should have %v errors, got %v", len(expectedErrors), len(result.Errors))
	}
	if result.IsValid != false {
		t.Fatalf("IsValid should be false, got %v", result.IsValid)
	}
	if !reflect.DeepEqual(expectedErrors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedErrors, result.Errors))
	}

}
func expectPassesRule(t *testing.T, rule graphql.ValidationRuleFn, queryString string) {
	expectValid(t, defaultSchema, []graphql.ValidationRuleFn{rule}, queryString)
}
func expectFailsRule(t *testing.T, rule graphql.ValidationRuleFn, queryString string, expectedErrors []gqlerrors.FormattedError) {
	expectInvalid(t, defaultSchema, []graphql.ValidationRuleFn{rule}, queryString, expectedErrors)
}

func ruleError(message string, locs ...int) gqlerrors.FormattedError {
	locations := []location.SourceLocation{}
	for i := 0; i < len(locs); i = i + 2 {
		line := locs[i]
		col := 0
		if i+1 < len(locs) {
			col = locs[i+1]
		}
		locations = append(locations, location.SourceLocation{
			Line:   line,
			Column: col,
		})
	}
	return gqlerrors.FormattedError{
		Message:   message,
		Locations: locations,
	}
}