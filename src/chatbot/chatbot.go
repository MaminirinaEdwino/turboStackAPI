package chatbot

import (
	"embed"
	_ "embed"
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

//go:embed model/*
var ModelDir embed.FS

var Reseaux *ReseauNeurones

// --- 1. Structures du Réseau de Neurones (MLP) ---
type Neurone struct {
	Poids    []float64 `json:"poids"`
	Biais    float64   `json:"biais"`
	Sortie   float64   `json:"sortie"`
	Gradient float64   `json:"gradient"`
}

type Couche struct {
	Neurones []Neurone `json:"neurones"`
	Entrees  []float64 `json:"entrees"`
}

type ReseauNeurones struct {
	Couches           []Couche `json:"couches"`
	TauxApprentissage float64  `json:"taux_apprentissage"`
}

// Structure globale pour sauvegarder l'état complet du modèle (Réseau + Vocabulaire)
type ModeleSauvegarde struct {
	Reseau      ReseauNeurones `json:"reseau"`
	MotsUniques []string       `json:"mots_uniques"`
	TagsIntents []string       `json:"tags_intents"`
}

// --- 2. Structures du Dataset (Intents) ---

type Intent struct {
	Tag       string   // Le nom de l'intention (e.g., "salutation")
	Patterns  []string // Les phrases d'entraînement
	Responses []string // Les réponses possibles
}

type TrainingPaire struct {
	Entrees          []float64 // Vecteur Bag-of-Words
	SortiesAttendues []float64 // Vecteur One-Hot
}

// --- 3. Fonctions Utilitaires (MLP) ---

func Sigmoide(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func DerivSigmoide(sortie float64) float64 {
	return sortie * (1.0 - sortie)
}

func AleatoirePoids(taille int) []float64 {
	poids := make([]float64, taille)
	for i := range poids {
		poids[i] = rand.Float64()*0.2 - 0.1
	}
	return poids
}

func CreerReseau(tailleEntree, tailleCache, tailleSortie int, taux float64) *ReseauNeurones {
	rn := &ReseauNeurones{TauxApprentissage: taux}

	// Couche Cachée
	coucheCachee := Couche{}
	for i := 0; i < tailleCache; i++ {
		coucheCachee.Neurones = append(coucheCachee.Neurones, Neurone{
			Poids: AleatoirePoids(tailleEntree),
			Biais: rand.Float64() * 0.1,
		})
	}
	rn.Couches = append(rn.Couches, coucheCachee)

	// Couche de Sortie
	coucheSortie := Couche{}
	for i := 0; i < tailleSortie; i++ {
		coucheSortie.Neurones = append(coucheSortie.Neurones, Neurone{
			Poids: AleatoirePoids(tailleCache),
			Biais: rand.Float64() * 0.1,
		})
	}
	rn.Couches = append(rn.Couches, coucheSortie)

	return rn
}

func (rn *ReseauNeurones) Propager(entrees []float64) []float64 {
	rand.Seed(time.Now().UnixNano())
	coucheCouranteEntrees := entrees

	for i := range rn.Couches {
		var sortiesCouche []float64
		rn.Couches[i].Entrees = coucheCouranteEntrees

		for j := range rn.Couches[i].Neurones {
			neurone := &rn.Couches[i].Neurones[j]
			sommePonderee := neurone.Biais
			for k, poids := range neurone.Poids {
				sommePonderee += poids * coucheCouranteEntrees[k]
			}
			neurone.Sortie = Sigmoide(sommePonderee)
			sortiesCouche = append(sortiesCouche, neurone.Sortie)
		}
		coucheCouranteEntrees = sortiesCouche
	}
	return coucheCouranteEntrees
}

func (rn *ReseauNeurones) Retropropager(sortiesAttendues []float64) {
	// Calculer les gradients
	for i := len(rn.Couches) - 1; i >= 0; i-- {
		couche := &rn.Couches[i]

		for j := range couche.Neurones {
			neurone := &couche.Neurones[j]
			erreur := 0.0

			if i == len(rn.Couches)-1 {
				erreur = sortiesAttendues[j] - neurone.Sortie
			} else {
				coucheSuivante := rn.Couches[i+1]
				for _, neuroneSuivant := range coucheSuivante.Neurones {
					erreur += neuroneSuivant.Gradient * neuroneSuivant.Poids[j]
				}
			}
			neurone.Gradient = erreur * DerivSigmoide(neurone.Sortie)
		}
	}

	// Mettre à jour les poids et les biais
	for i := range rn.Couches {
		couche := &rn.Couches[i]
		entrees := couche.Entrees

		for j := range couche.Neurones {
			neurone := &couche.Neurones[j]

			for k := range neurone.Poids {
				neurone.Poids[k] += rn.TauxApprentissage * neurone.Gradient * entrees[k]
			}
			neurone.Biais += rn.TauxApprentissage * neurone.Gradient
		}
	}
}

// --- 4. PERSISTENCE (Sauvegarde & Chargement) ---



func ChargerModele() (*ReseauNeurones, []string, []string, error) {
	fichier, err := ModelDir.Open("model/model_chatbot.json")
	if err != nil {
		return nil, nil, nil, err
	}
	defer fichier.Close()

	var modele ModeleSauvegarde
	decodeur := json.NewDecoder(fichier)
	if err := decodeur.Decode(&modele); err != nil {
		return nil, nil, nil, err
	}
	

	return &modele.Reseau, modele.MotsUniques, modele.TagsIntents, nil
}

// --- 6. PRÉTRAITEMENT (Vectorisation) ---

var (
	MotsUniques []string // Le vocabulaire entier (Bag-of-Words)
	TagsIntents []string // Ordre des tags pour le One-Hot Encoding
)

func preparerDataset(intents []Intent) (trainData []TrainingPaire) {
	motsTemp := make(map[string]bool)

	for _, intent := range intents {
		TagsIntents = append(TagsIntents, intent.Tag)
		for _, pattern := range intent.Patterns {
			tokens := strings.FieldsSeq(strings.ToLower(pattern))
			for token := range tokens {
				motsTemp[token] = true
			}
		}
	}

	for mot := range motsTemp {
		MotsUniques = append(MotsUniques, mot)
	}
	sort.Strings(MotsUniques)

	for i, intent := range intents {
		output := make([]float64, len(TagsIntents))
		output[i] = 1.0

		for _, pattern := range intent.Patterns {
			input := make([]float64, len(MotsUniques))
			tokens := strings.FieldsSeq(strings.ToLower(pattern))

			for token := range tokens {
				for j, mot := range MotsUniques {
					if mot == token {
						input[j] = 1.0
						break
					}
				}
			}

			trainData = append(trainData, TrainingPaire{
				Entrees:          input,
				SortiesAttendues: output,
			})
		}
	}
	return trainData
}

func textToBoW(texte string) []float64 {
	bow := make([]float64, len(MotsUniques))
	tokens := strings.Fields(strings.ToLower(texte))

	for _, token := range tokens {
		for i, mot := range MotsUniques {
			if mot == token {
				bow[i] = 1.0
				break
			}
		}
	}
	return bow
}

// --- 7. INFÉRENCE (Chatbot Core) ---

func ObtenirIntention(rn *ReseauNeurones, phrase string, intents []Intent) string {
	rand.Seed(time.Now().UnixNano())
	inputVector := textToBoW(phrase)
	probabilities := rn.Propager(inputVector)

	maxIndex := -1
	maxProb := -1.0
	for i, prob := range probabilities {
		if prob > maxProb {
			maxProb = prob
			maxIndex = i
		}
	}

	if maxProb < 0.75 {
		return "Je ne suis pas sûr d'avoir compris. Pouvez-vous reformuler ?"
	}

	predictedTag := TagsIntents[maxIndex]

	for _, intent := range intents {
		if intent.Tag == predictedTag {
			return intent.Responses[rand.Intn(len(intent.Responses))]
		}
	}

	return "Je ne sais pas quoi répondre à cela."
}

// --- 8. MAIN et Chargement / Entraînement ---

// func chatmain() {
// 	rand.Seed(time.Now().UnixNano())

// 	intents := ChargerIntents()
// 	rn, _, _, err := ChargerModele()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Tentative de chargement du modèle existant

// 	// Si le modèle n'a pas pu être chargé, on l'entraîne de zéro

// 	// --- BOUCLE D'INTERACTION ---
// 	fmt.Println("\n--- Chatbot TurboStack Prêt (tapez 'quitter' pour sortir) ---")
// 	scanner := bufio.NewScanner(os.Stdin)

// 	for {
// 		fmt.Print("Vous: ")
// 		if !scanner.Scan() {
// 			break
// 		}
// 		input := scanner.Text()

// 		if strings.ToLower(strings.TrimSpace(input)) == "quitter" {
// 			fmt.Println("Bot: Au revoir !")
// 			break
// 		}

// 		reponse := ObtenirIntention(rn, input, intents)
// 		fmt.Printf("Bot: %s\n", reponse)
// 	}
// }
