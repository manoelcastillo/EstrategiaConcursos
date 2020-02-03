package handler

import (
	"EstrategiaConcursos/model"
	"strings"
)

//a ideia do algoritmo é sugerir tags baseando-se na descrição do repositório e também na linguagem. Cpm certeza cabe vários "tweaks", mas para o exemplo serve.
//possiveis comportamentos indesejados: na descrição pode haver tags repetidas, podem aparecer tags sem sentido como "therefore, there, that", enfim, palavras sem
//nexo para serem sugeridas como topic do repositório.
//Possiveis melhorias: 1. Ver se as palavras filtradas estão presentes em uma lista pre-definida de topics 2. Uso de machine learning.
func GetPossibleTags(repo model.Repository) []model.Tag {
	split := strings.Fields(repo.Lang + " " + repo.Desc)
	var tags []model.Tag
	
	for _, x := range split {
		var t model.Tag
		t.Name = x
		if (len(x) > 3 || strings.EqualFold(x, "api") || strings.EqualFold(x, "xml")) && !TagInUse(repo, t) {
			tags = append(tags, t)
		}
	}

	return tags
}

func TagInUse(repo model.Repository, tag model.Tag) bool {
	for _, x := range repo.Tags {
		if strings.EqualFold(x.Name, tag.Name) {
			return true
		}
	}
	return false
}
