package chatbot

func ChargerIntents() []Intent {
	return []Intent{
		// 1. SALUTATION
		{
			Tag: "salutation",
			Patterns: []string{
				"Bonjour", "Salut", "Coucou", "Hello", "Hey", "Hi",
				"Comment ça va?", "Ça va?", "Yop", "Bonjour ça va?", 
				"J'ai besoin d'aide", "Salut l'ami", "Bonsoir", "Bonne journée",
			},
			Responses: []string{
				"Bonjour ! Comment puis-je vous aider aujourd'hui ?",
				"Salut ! Quel est le sujet de votre visite ?",
				"Hello ! Je suis à votre écoute pour toute question.",
				"Bonsoir ! Que puis-je faire pour vous ?",
			},
		},
		
		// 2. REMERCIEMENTS
		{
			Tag: "remerciement",
			Patterns: []string{
				"Merci", "Merci beaucoup", "Je vous remercie", 
				"C'est gentil", "Super, merci", "Merci pour l'aide",
				"Génial merci", "Top, merci à vous", "Cool merci",
			},
			Responses: []string{
				"Avec plaisir ! N'hésitez pas si vous avez d'autres questions.",
				"De rien ! C'est un plaisir de vous aider.",
				"Je vous en prie. Bonne continuation !",
				"Content d'avoir pu vous aider !",
			},
		},
		
		// 3. COMMANDE DE PRODUIT / PIZZA
		{
			Tag: "commande_produit",
			Patterns: []string{
				"Je veux commander", "Commander un article", 
				"Je commande une pizza", "Une pizza svp", 
				"J'aimerais acheter une quatre saisons", "Je veux passer commande",
				"Je voudrais commander à manger", "Puis-je commander ?",
				"Je souhaite acheter un produit",
			},
			Responses: []string{
				"Quelle taille ou quelle option souhaitez-vous ? (Petite, moyenne, grande)",
				"D'accord, quel type d'article ou de pizza désirez-vous ?",
				"C'est parti pour la commande ! Donnez-moi plus de détails.",
			},
		},
		
		// 4. STATUT DE COMMANDE / LIVRAISON
		{
			Tag: "statut_commande",
			Patterns: []string{
				"Où est ma commande ?", "Statut de ma livraison", 
				"Quand ma commande arrive-t-elle ?", "Ma commande est-elle prête ?", 
				"Suivre mon colis", "Le livreur est-il en route ?",
				"Où en est mon achat ?", "Temps de livraison",
			},
			Responses: []string{
				"Pouvez-vous me donner votre numéro de commande pour que je vérifie ?",
				"Je consulte le suivi. Quel est votre nom de famille ?",
				"La livraison prend généralement entre 30 et 45 minutes.",
			},
		},

		// 5. LOCALISATION ET HORAIRES
		{
			Tag: "localisation_horaires",
			Patterns: []string{
				"Où êtes-vous situés ?", "Quelle est l'adresse ?", 
				"Quel est le magasin le plus proche ?", "Vos heures d'ouverture",
				"Êtes-vous ouverts le dimanche ?", "Quels sont vos horaires",
				"Adresse de la boutique", "C'est où ?",
			},
			Responses: []string{
				"Nous sommes situés au 123 Rue Principale. Ouvert de 11h à 23h tous les jours.",
				"Notre établissement est ouvert du Lundi au Dimanche de 11h00 à 23h00.",
				"Vous pouvez nous rendre visite à notre adresse principale ou commander en ligne !",
			},
		},

		// 6. MOYENS DE PAIEMENT
		{
			Tag: "paiement",
			Patterns: []string{
				"Quels sont les modes de paiement ?", "Acceptez-vous la carte bancaire ?", 
				"Puis-je payer en espèces ?", "Paiement mobile disponible ?", 
				"Comment régler ma commande ?", "Cartes acceptées",
			},
			Responses: []string{
				"Nous acceptons les cartes bancaires, les paiements mobiles et les espèces à la livraison.",
				"Vous pouvez régler en ligne par carte ou directement à la livraison.",
			},
		},

		// 7. RÉCLAMATION / PROBLÈME
		{
			Tag: "reclamation",
			Patterns: []string{
				"J'ai un problème", "Ma commande est incomplète", 
				"Le produit est endommagé", "Je veux un remboursement", 
				"Il y a une erreur dans ma commande", "Insatisfait",
			},
			Responses: []string{
				"Je suis désolé pour ce désagrément. Pouvez-vous me donner plus de détails sur le problème ?",
				"Nous allons régler cela rapidement. Quel est votre numéro de commande ?",
			},
		},
		
		// 8. ADIEU
		{
			Tag: "adieu",
			Patterns: []string{
				"Au revoir", "À plus", "Ciao", "Merci et au revoir", 
				"Je dois partir", "Bonne journée", "Bye bye", "À bientot",
			},
			Responses: []string{
				"Au revoir ! À très bientôt.",
				"Merci de votre visite et excellente journée à vous !",
				"Prenez soin de vous, à la prochaine !",
			},
		},
		// 9. CRÉATION DE PROJET TURBOSTACK
		{
			Tag: "creation_projet",
			Patterns: []string{
				"Comment créer un nouveau projet ?", "Je veux créer une application", 
				"Comment démarrer un projet TurboStack ?", "Lancer un nouveau projet",
				"Créer une API ou un site", "Nouveau projet",
			},
			Responses: []string{
				"Pour créer un projet, rendez-vous sur le tableau de bord et cliquez sur le bouton de création pour définir votre nom et type de projet.",
				"TurboStack vous permet de lancer rapidement un projet en quelques clics, sans code complexe.",
			},
		},

		// 10. TYPES DE PROJETS (API, WebApp, Site Statique)
		{
			Tag: "types_projets",
			Patterns: []string{
				"Quels types de projets puis-je faire ?", "C'est quoi une WebApp ?", 
				"Quelle différence entre site statique et API ?", "Types d'applications supportés",
				"Faire une API REST", "Créer un site web statique",
			},
			Responses: []string{
				"TurboStack prend en charge trois types principaux : les API REST, les applications web interactives (WebApp) et les sites web statiques.",
				"Vous pouvez choisir de structurer une API backend, une interface web dynamique ou un site vitrine selon vos besoins.",
			},
		},

		// 11. BASE DE DONNÉES ET TABLES
		{
			Tag: "base_de_donnees",
			Patterns: []string{
				"Comment gérer la base de données ?", "Créer une table", 
				"C'est quoi un Table node ?", "Comment lier des tables ?", 
				"Utiliser PostgreSQL avec TurboStack", "Query builder c'est quoi ?",
			},
			Responses: []string{
				"Le module de base de données vous permet de créer de nouvelles tables, d'ajouter des champs et de structurer vos données visuellement.",
				"Les 'Table nodes' affichent vos tables sous forme de cartes interactives, et le Query builder vous aide à concevoir vos requêtes facilement.",
			},
		},

		// 12. EXPORTATION DE PROJET
		{
			Tag: "exportation_projet",
			Patterns: []string{
				"Comment exporter mon projet ?", "Où sont les fichiers générés ?", 
				"Télécharger le code source", "Récupérer mon application",
				"Exportation Go et code source",
			},
			Responses: []string{
				"L'option d'exportation rassemble tout votre travail et génère proprement les fichiers de votre projet (notamment en Go et SQL) prêts à être exécutés.",
				"Vous pouvez exporter l'intégralité de votre architecture pour l'exécuter de manière totalement autonome.",
			},
		},
		{
			Tag: "types_projets",
			Patterns: []string{
				"Quels types de projets puis-je créer ?", 
				"Quels sont les choix de projets disponibles ?", 
				"Différence entre site statique et application web", 
				"Créer une API REST avec TurboStack", 
				"Structure d'un projet webapp",
			},
			Responses: []string{
				"TurboStack gère trois grands types de projets : les API REST (back-end), les applications web dynamiques (WebApp) et les sites web statiques.",
				"Chaque type de projet initialise une architecture de dossiers dédiée (comme /api, /views, /static) prête pour le développement.",
			},
		},

		// 2. CONFIGURATION DE LA BASE DE DONNÉES ET MODÈLES
		{
			Tag: "bdd_et_modeles",
			Patterns: []string{
				"Comment configurer la base de données ?", 
				"Créer un modèle de données", 
				"Ajouter des tables et des champs", 
				"Gérer les types de champs et contraintes", 
				"Initialisation de PostgreSQL avec TurboStack",
			},
			Responses: []string{
				"L'éditeur de base de données permet de définir vos modèles, leurs attributs (entiers, chaînes de caractères) et leurs contraintes (clés primaires, relations).",
				"TurboStack génère automatiquement les scripts d'initialisation SQL (tables et requêtes CRUD) pour PostgreSQL à partir de vos modèles.",
			},
		},

		// 3. ÉDITEUR D'API ET ENDPOINTS
		{
			Tag: "editeur_api",
			Patterns: []string{
				"Comment configurer les endpoints de l'API ?", 
				"Créer des routes HTTP POST GET PUT DELETE", 
				"Associer un modèle à un endpoint", 
				"Gérer les paramètres d'URL et du Body", 
				"Génération des contrôleurs Go",
			},
			Responses: []string{
				"L'éditeur d'API permet de lier des méthodes HTTP (GET, POST, etc.) à des URI spécifiques et de définir les contenus de retour ou de redirection.",
				"Le moteur génère le code source Go correspondant dans le dossier des contrôleurs et met à jour le routeur `http.ServeMux` automatiquement.",
			},
		},

		// 4. INTERFACES VISUELLES ET PAGES WEB
		{
			Tag: "pages_et_frontend",
			Patterns: []string{
				"Comment concevoir les pages web ?", 
				"Utiliser le constructeur visuel de pages", 
				"Ajouter des composants et des blocs HTML", 
				"Gérer les styles CSS par média query", 
				"Lier les données de l'API aux éléments de la page",
			},
			Responses: []string{
				"Le constructeur de pages utilise une hiérarchie de blocs (conteneurs, formulaires, champs de saisie, images) modifiables visuellement.",
				"Les styles (bureau, tablette, mobile) sont automatiquement compilés et associés aux fichiers HTML de sortie de votre application.",
			},
		},

		// 5. EXPORTATION ET BUILD DU PROJET
		{
			Tag: "exportation_projet",
			Patterns: []string{
				"Comment exporter un projet complet ?", 
				"Génération des fichiers sources Go et HTML", 
				"Où trouver les fichiers exportés", 
				"Compiler le projet généré", 
				"Structure du code exporté par TurboStack",
			},
			Responses: []string{
				"L'exportation génère l'arborescence complète du projet incluant le fichier `main.go`, les configurations de base de données, les routes et les vues HTML.",
				"Vous récupérez un projet autonome prêt à être compilé et exécuté avec l'environnement Go et PostgreSQL.",
			},
		},

		// 6. INTERFACE D'ADMINISTRATION GLOBALE
		{
			Tag: "admin_dashboard",
			Patterns: []string{
				"Comment fonctionne l'interface d'administration ?", 
				"Voir la liste de tous les projets utilisateurs", 
				"Gérer les utilisateurs depuis l'admin", 
				"Supprimer un projet par son nom", 
				"Authentification admin par JWT",
			},
			Responses: []string{
				"L'interface d'administration sécurisée par jeton JWT permet de lister l'ensemble des utilisateurs inscrits et tous leurs projets globaux.",
				"Le panneau d'administration offre la possibilité d'examiner en détail les structures de projets ou d'effectuer des suppressions ciblées.",
			},
		},
	}
}