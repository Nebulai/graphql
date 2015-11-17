package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_NoCircularFragmentSpreads_SingleReferenceIsValid(t *testing.T) {
	expectPassesRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB }
      fragment fragB on Dog { name }
    `)
}
func TestValidate_NoCircularFragmentSpreads_SpreadingTwiceIsNotCircular(t *testing.T) {
	expectPassesRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB, ...fragB }
      fragment fragB on Dog { name }
    `)
}
func TestValidate_NoCircularFragmentSpreads_SpreadingTwiceIndirectlyIsNotCircular(t *testing.T) {
	expectPassesRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB, ...fragC }
      fragment fragB on Dog { ...fragC }
      fragment fragC on Dog { name }
    `)
}
func TestValidate_NoCircularFragmentSpreads_DoubleSpreadWithinAbstractTypes(t *testing.T) {
	expectPassesRule(t, graphql.NoFragmentCyclesRule, `
      fragment nameFragment on Pet {
        ... on Dog { name }
        ... on Cat { name }
      }

      fragment spreadsInAnon on Pet {
        ... on Dog { ...nameFragment }
        ... on Cat { ...nameFragment }
      }
    `)
}
func TestValidate_NoCircularFragmentSpreads_SpreadingRecursivelyWithinFieldFails(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Human { relatives { ...fragA } },
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself.`, 2, 45),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDirectly(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragA }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself.`, 2, 31),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDirectlyWithinInlineFragment(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Pet {
        ... on Dog {
          ...fragA
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself.`, 4, 11),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDirectlyMultiple(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB }
      fragment fragB on Dog { ...fragA }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself via fragB.`, 2, 31, 3, 31),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDirectlyWithinInlineFragmentMultiple(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Pet {
        ... on Dog {
          ...fragB
        }
      }
      fragment fragB on Pet {
        ... on Dog {
          ...fragA
        }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself via fragB.`, 4, 11, 9, 11),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDeeply(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB }
      fragment fragB on Dog { ...fragC }
      fragment fragC on Dog { ...fragO }
      fragment fragX on Dog { ...fragY }
      fragment fragY on Dog { ...fragZ }
      fragment fragZ on Dog { ...fragO }
      fragment fragO on Dog { ...fragA, ...fragX }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself via fragB, fragC, fragO.`, 2, 31, 3, 31, 4, 31, 8, 31),
		ruleError(`Cannot spread fragment "fragX" within itself via fragY, fragZ, fragO.`, 5, 31, 6, 31, 7, 31, 8, 41),
	})
}
func TestValidate_NoCircularFragmentSpreads_NoSpreadingItselfDeeplyTwoPaths(t *testing.T) {
	expectFailsRule(t, graphql.NoFragmentCyclesRule, `
      fragment fragA on Dog { ...fragB, ...fragC }
      fragment fragB on Dog { ...fragA }
      fragment fragC on Dog { ...fragA }
    `, []gqlerrors.FormattedError{
		ruleError(`Cannot spread fragment "fragA" within itself via fragB.`, 2, 31, 3, 31),
		ruleError(`Cannot spread fragment "fragA" within itself via fragC.`, 2, 41, 4, 31),
	})
}