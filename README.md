# Loganalyzer – Analyse de logs concurrente (TP Go)
 
Bienvenue dans ce petit projet Go ! On y développe `loganalyzer`, un outil en ligne de commande pensé pour illustrer la lecture d'un fichier JSON, l'exécution concurrente avec des goroutines et la gestion propre des erreurs. Tout est volontairement simple pour rester pédagogique et aller droit au but.
 
## ✨ Fonctionnalités principales
- **Commande `analyze`** : lit un fichier de configuration JSON et lance une goroutine par log à inspecter.
- **Analyse simulée** : chaque fichier est ouvert pour vérifier son accessibilité, puis on simule un traitement (50 à 200 ms) avec **10% de chance de lever une erreur d'analyse**.
- **Résultats centralisés** : affichage en console (avec un résumé chaleureux) et export possible vers un rapport JSON.
- **Erreurs personnalisées** : différenciation claire entre les soucis d'accès aux fichiers et les erreurs de parsing.
## 🚀 Prise en main rapide
 
----
```bash
git clone <votre_repo_git> loganalyzer
cd loganalyzer
go run . analyze --config config.json
```
 
Astuce : `go install` permet d'installer le binaire dans votre `$GOBIN`.
 
## 🧾 Format du fichier de configuration
 
Le fichier attendu (`--config` ou `-c`) est un tableau JSON :
 
```json
[
  { "id": "web-server-1", "path": "test_logs/access.log", "type": "nginx-access" },
  { "id": "app-backend-2", "path": "test_logs/errors.log", "type": "custom-app" }
]
```
 
Chaque entrée devient une tâche d'analyse concurrente.

## 📦 Export d'un rapport JSON
 
Pour sauvegarder les résultats, ajoutez `--output` (ou `-o`) :
 
```bash
go run . analyze -c config.json -o report.json
```
 
Le rapport contient pour chaque log :
 
```json
[
  {
    "log_id": "web-server-1",
    "file_path": "test_logs/access.log",
    "status": "OK",
    "message": "Analyse terminée avec succès.",
    "error_details": ""
  },
  {
    "log_id": "invalid-path",
    "file_path": "/non/existent/log.log",
    "status": "FAILED",
    "message": "Fichier introuvable ou illisible.",
    "error_details": "invalid-path (/non/existent/log.log): open /non/existent/log.log: no such file or directory"
  }
]
```

Si le chemin cible inclut des répertoires inexistants, ils sont créés automatiquement.

## 🖨️ Exemple de sortie console

---
Une exécution avec le fichier `config.json` du dépôt donne quelque chose comme ceci :

```
✅ web-server-1 — Analyse réussie (test_logs/access.log)
   ↳ Analyse terminée avec succès.
❌ db-server-3 — Analyse en panne (test_logs/mysql_error.log)
   ↳ Fichier introuvable ou illisible.
   ⚠️ Détail: open test_logs/mysql_error.log: no such file or directory
✅ app-backend-2 — Analyse réussie (test_logs/errors.log)
   ↳ Analyse terminée avec succès.

✨ Bilan: 2 réussite(s), 1 échec(s) sur 3 analyse(s).
```

Selon les 10% de chance d'échec simulé, un log peut parfois remonter une erreur d'analyse supplémentaire. Dans ce cas, un bloc `⚠️ étail` décrit la raison du faux pas.

---
## ⚙️ Gestion des erreurs

Deux types d'erreurs personnalisées sont exposés et manipulés via `errors.Is` / `errors.As` :

| Cas | Type | Description |
| --- | --- | --- |
| Accès impossible au fichier de configuration | `config.FileError` | Associe le chemin incriminé à l'erreur système.
| JSON invalide | `config.ParseError` | Remonte le détail du parsing.
| Fichier de log illisible | `analyzer.LogFileError` | Fournit l'identifiant du log concerné.

Ces erreurs sont capturées pour afficher des messages clairs et compréhensibles.

## 🧱 Architecture du code

```
.
├── cmd/               # Commandes Cobra (root + analyze)
├── internal/
│   ├── analyzer/      # Lancement concurrent et structuration des résultats
│   ├── config/        # Lecture & validation du JSON de configuration
│   └── reporter/      # Export des rapports JSON
├── test_logs/         # Quelques fichiers de tests
└── main.go            # Point d'entrée du programme
```

Chaque package exporte des structures et fonctions documentées pour faciliter la lecture.

## 🧪 Tests / Vérifications

Le projet ne contient pas de tests unitaires formels, mais la commande suivante vérifie la compilation de l'ensemble :

```bash
go test ./...
```

> Pensez à lancer également `go run . analyze -c config.json` pour voir le fonctionnement réel.


<img width="1046" height="341" alt="image" src="https://github.com/user-attachments/assets/28415d42-164c-4080-955c-5d76c8076ce1" />

