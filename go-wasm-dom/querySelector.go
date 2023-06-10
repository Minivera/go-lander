package go_wasm_dom

import (
	"fmt"
	"math"
)

type tokenType uint64

const (
	tag tokenType = iota
	attribute
	id
	class
)

type selector struct {
	selector     string
	selectorType tokenType
	negated      bool
}

type chain struct {
	isCurrent                bool
	shouldBeDirectDescendant bool
	selectors                []selector
}

type token struct {
	content  string
	posStart int
}

func selectorsFromToken(token token) ([]selector, error) {
	// Loop over all the characters of the token one by one. We are searching for any special characters
	// and creating new selectors based on that
	var selectors []selector
	var currentSelector *selector
	checkingForPseudo := false
	shouldNegate := false
	pseudoSelector := ""
	for pos, char := range token.content {
		switch char {
		case '.':
			// If hitting a dot, this is a class selector. Append the previous selector if it existed
			// and create a new one.
			if currentSelector != nil {
				selectors = append(selectors, *currentSelector)
			}

			currentSelector = &selector{
				selectorType: class,
				negated:      shouldNegate,
			}
		case '#':
			// Same for the number character, which defines an id selector
			if currentSelector != nil {
				selectors = append(selectors, *currentSelector)
			}

			currentSelector = &selector{
				selectorType: id,
				negated:      shouldNegate,
			}
		case '[':
			// Again same for the opening square bracket, that defines an attribute selector.
			if currentSelector != nil {
				selectors = append(selectors, *currentSelector)
			}

			currentSelector = &selector{
				selectorType: attribute,
				negated:      shouldNegate,
			}
		case ']':
			// If we hit a closing square bracket, then we're losing a selector. We should make sure we
			// opened a valid attribute selector or return an error.
			if currentSelector != nil && currentSelector.selectorType == attribute {
				selectors = append(selectors, *currentSelector)
			} else {
				return nil, fmt.Errorf(
					"selector %s is not valid, character ']' at position %d was used without an opening '['",
					token.content,
					token.posStart+pos,
				)
			}
		case ':':
			// If we hit a colon, then we're processing a pseudo selector. The next characters processed should
			// be ignored until we hit an opening parenthesis.
			checkingForPseudo = true
			pseudoSelector = ""
		case '(':
			// If we hit an opening parenthesis, the pseudo selector is done and we can act on it
			if !checkingForPseudo {
				return nil, fmt.Errorf(
					"selector %s is not valid, character '(' at position %d was used without a pseudo selector",
					token.content,
					token.posStart+pos,
				)
			}

			switch pseudoSelector {
			case "not":
				checkingForPseudo = false
				shouldNegate = true
			default:
				return nil, fmt.Errorf(
					"pseudo selector :%s is either unknown or invalid",
					pseudoSelector,
				)
			}
		case ')':
			// If we hit a closing parenthesis, then we are closing whatever pseudo selector we opened. Process
			// it.
			shouldNegate = false
		default:
			// Everything else, we either add to the current selector if it exists, thread as a tag if not
			// or treat as a pseudo selector
			if checkingForPseudo {
				pseudoSelector += string(char)
				continue
			}

			if currentSelector == nil {
				currentSelector = &selector{
					selectorType: tag,
					selector:     string(char),
					negated:      shouldNegate,
				}
				continue
			}

			currentSelector.selector += string(char)
		}
	}

	return selectors, nil
}

func tokenize(selectors string) ([]chain, error) {
	// Tokenize the string by extracting any terms between spaces. We don't
	// use split to keep the position of the tokens for error messages.
	var selectorParts []token
	var currentPart *token
	for pos, char := range selectors {
		if char != ' ' {
			if currentPart == nil {
				currentPart = &token{
					content:  string(char),
					posStart: pos,
				}
				continue
			}

			currentPart.content += string(char)
		} else if char == ' ' && currentPart != nil {
			selectorParts = append(selectorParts, *currentPart)
			currentPart = nil
		}
	}

	// Then take each token and generate all its selectors into a chain.
	// A token can be something like input.test[data-test="value"], we want to extract the relevant
	// parts of this selector into something we can process.
	var selectorChain []chain
	var currentSelector *chain
	for _, part := range selectorParts {
		// This checks for direct descendants, update the current chain to set the property and look for
		// the next selector.
		if part.content == ">" {
			if currentSelector == nil {
				return nil, fmt.Errorf(
					"selector %s is not valid, character '>' at position %d does not follow a parent selector. You might have forgotten to add '&'",
					selectors,
					currentPart.posStart,
				)
			}
			currentSelector.shouldBeDirectDescendant = true
			continue
		}

		// This checks for the "current" element. It should only ever be the first selector passed
		// and no property should be set on that selector.
		if part.content == "&" {
			if currentSelector != nil {
				return nil, fmt.Errorf(
					"selector %s is not valid, character '&' at position %d cannot be used at this position",
					selectors,
					currentPart.posStart,
				)
			}

			currentSelector = &chain{
				isCurrent: true,
			}
			continue
		}

		if currentSelector != nil {
			selectorChain = append(selectorChain, *currentSelector)
		}

		selectors, err := selectorsFromToken(part)
		if err != nil {
			return nil, err
		}

		currentSelector.selectors = selectors
	}

	return selectorChain, nil
}

func selectAll(current Value, selectors string) []Value {
	chains, err := tokenize(selectors)
	if err != nil {
		t.Fatalf("tried to execute selector, but the selector could not be parsed. %v", err)
	}

	var candidates []int
	for _, chain := range chains {
		// If we're checking the current parent, then get all nodes where the given parent
		// is either a parent or descendant
		if chain.isCurrent {
			if chain.shouldBeDirectDescendant {
				// Use the descendant map
				for candidate, ancestors := range currentScreen.nodePerAncestors {
					// The descendants are a slice, so loop over all descendants
					for _, ancestor := range ancestors {
						// Once we've found a descendant, we can stop
						if ancestor == current.id {
							candidates = append(candidates, candidate)
							continue
						}
					}
				}
				continue
			}
			// Otherwise, use the parents map
			for candidate, parent := range currentScreen.nodePerParent {
				if parent == current.id {
					candidates = append(candidates, candidate)
				}
			}
			continue
		}

		// Get the temporary candidates based on the selector
		// TODO: Implement the negations
		var tempCandidates []int
		for _, selector := range chain.selectors {
			// Get all nodes that match
			var subTempCandidates []int
			switch selector.selectorType {
			case id:
				obtained, ok := currentScreen.nodesPerID[selector.selector]
				if ok {
					subTempCandidates = obtained
				}
			case class:
				obtained, ok := currentScreen.nodesPerClass[selector.selector]
				if ok {
					subTempCandidates = obtained
				}
			case tag:
				obtained, ok := currentScreen.nodesPerTag[selector.selector]
				if ok {
					subTempCandidates = obtained
				}
			case attribute:
				obtained, ok := currentScreen.nodesPerAttribute[selector.selector]
				if ok {
					subTempCandidates = obtained
				}
			}

			// Filter out any nodes that is not in the temp candidates or was not found here
			actualCandidates := make([]int, 0, int(math.Max(
				float64(len(tempCandidates)),
				float64(len(subTempCandidates)),
			)))
			// To be kept, a candidate must match the previous selectors, and the current selector
		subTempLoop:
			for _, c := range subTempCandidates {
				for _, c2 := range tempCandidates {
					if c == c2 {
						actualCandidates = append(actualCandidates, c)
						continue subTempLoop
					}
				}
			}
			tempCandidates = actualCandidates
		}

		// Filter out the temporary candidates based on the previous candidates, if any
		actualCandidates := make([]int, 0, len(tempCandidates))
	tempLoop:
		for _, temp := range tempCandidates {
			// To be kept, a candidate must be either a descendant of another candidate, or a direct
			// parent if that's the current setting
			if chain.shouldBeDirectDescendant {
				for _, parent := range candidates {
					// Check if any of the previous candidates is a parent of the current temporary candidate
					if p, ok := currentScreen.nodePerParent[temp]; ok && p == parent {
						// If yet, keep it
						actualCandidates = append(actualCandidates, temp)
						continue tempLoop
					}
				}
			} else {
				for _, ancestor := range candidates {
					// Check if any of the previous candidates is an ancestor of the current temporary candidate
					ancestors, ok := currentScreen.nodePerAncestors[temp]
					if !ok {
						continue
					}

					// Check all the ancestors of this node in order
					for _, a := range ancestors {
						if a == ancestor {
							// If any ancestor match, keep the candidate
							actualCandidates = append(actualCandidates, temp)
							continue tempLoop
						}
					}
				}
			}
		}

		candidates = tempCandidates
	}

	var found []Value
candidateLoop:
	for _, candidate := range candidates {
		for _, v := range currentScreen.allNodes {
			if v.id == candidate {
				found = append(found, *v)
				continue candidateLoop
			}
		}
	}

	return found
}

func querySelectorAll(caller Value, args ...Value) []Value {
	if currentScreen == nil {
		t.Fatal("querySelectorAll tried to execute, but tests were not initialized")
	}

	if caller.referencedType != domNode {
		t.Fatalf("querySelectorAll expects this to be a DOM node, %T given", caller)
	}

	if len(args) != 1 {
		t.Fatalf("querySelectorAll expects 1 argument, %d given", len(args))
	}

	selectors := args[0]
	if selectors.referencedType != stringType {
		t.Fatalf("querySelectorAll expects argument[0] to be a string, %T given", args[0])
	}

	return selectAll(caller, args[0].String())
}

func querySelector(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("querySelector tried to execute, but tests were not initialized")
	}

	if caller.referencedType != domNode {
		t.Fatalf("querySelector expects this to be a DOM node, %T given", caller)
	}

	if len(args) != 1 {
		t.Fatalf("querySelector expects 1 argument, %d given", len(args))
	}

	selectors := args[0]
	if selectors.referencedType != stringType {
		t.Fatalf("querySelector expects argument[0] to be a string, %T given", args[0])
	}

	found := selectAll(caller, args[0].String())
	if len(found) > 1 {
		t.Fatalf("querySelector expects to find only one node, but %d were found. Use querySelectorAll instead", found)
	}

	return found[0]
}
