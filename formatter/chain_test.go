package formatter

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/octoberswimmer/apexfmt/parser"

	log "github.com/sirupsen/logrus"
)

func TestChain(t *testing.T) {
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)

	}
	tests :=
		[]struct {
			input  string
			output int
		}{
			{`Schema.SObjectType.Account.getRecordTypeInfosByDeveloperName().get('Business').getRecordTypeId()`, 3},
			{`Fixtures.Contact(account).put(Contact.RecordTypeId, Schema.SObjectType.Contact.getRecordTypeInfosByDeveloperName().get('Person').getRecordTypeId()).put(Contact.My_Lookup__c, newRecord[0].Id).save()`, 4},
		}

	for _, tt := range tests {
		input := antlr.NewInputStream(tt.input)
		lexer := parser.NewApexLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		p := parser.NewApexParser(stream)
		p.RemoveErrorListeners()

		v := NewChainVisitor()
		out, ok := v.visitRule(p.Statement()).(int)
		if !ok {
			t.Errorf("Unexpected result parsing apex")
		}
		if out != tt.output {
			t.Errorf("unexpected result.  expected: %d; got: %d", tt.output, out)
		}
	}
}