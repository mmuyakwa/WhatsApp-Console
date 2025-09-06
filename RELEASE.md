# Release Anleitung - WhatsApp Console

## 🎯 Übersicht

Dieses Dokument beschreibt, wie du neue Releases für WhatsApp Console erstellst. Der Prozess ist durch GitHub Actions vollständig automatisiert.

## 📋 Voraussetzungen

- Push-Rechte auf das Repository
- Alle Änderungen sind auf dem `master` Branch
- Alle Tests laufen erfolgreich

## 🚀 Release Prozess

### Schritt 1: Version vorbereiten

```bash
# 1. Lokales Repository aktualisieren
git checkout master
git pull origin master

# 2. Alle Änderungen committen und pushen
git add .
git commit -m "Prepare release v1.2.0"
git push origin master

# 3. Workflows überprüfen
# Gehe zu: https://github.com/mmuyakwa/WhatsApp-Console/actions
# Stelle sicher, dass alle Workflows erfolgreich sind
```

### Schritt 2: Git Tag erstellen

```bash
# Semantic Versioning verwenden
# Format: v{MAJOR}.{MINOR}.{PATCH}
# 
# MAJOR: Breaking changes (v1.0.0 → v2.0.0)
# MINOR: Neue Features (v1.0.0 → v1.1.0)  
# PATCH: Bugfixes (v1.0.0 → v1.0.1)

# Tag mit aussagekräftiger Message erstellen
git tag -a v1.2.0 -m "Release v1.2.0: Add message scheduling and improved error handling"

# Tag zu GitHub pushen (startet automatisch Release-Workflow!)
git push origin v1.2.0
```

### Schritt 3: Release-Workflow überwachen

Der `release.yml` Workflow wird automatisch gestartet und:

1. **Führt Tests aus** - Stellt sicher, dass alles funktioniert
2. **Erstellt Cross-Platform Binaries**:
   - `whatsapp-console-linux-amd64` (Linux 64-bit)
   - `whatsapp-console-linux-arm64` (Linux ARM64)
   - `whatsapp-console-windows-amd64.exe` (Windows 64-bit)
   - `whatsapp-console-darwin-amd64` (macOS Intel)
   - `whatsapp-console-darwin-arm64` (macOS Apple Silicon)
3. **Erstellt GitHub Release** mit Download-Links
4. **Fügt automatisch Changelog hinzu**

### Schritt 4: Release verifizieren

1. Gehe zu [GitHub Releases](https://github.com/mmuyakwa/WhatsApp-Console/releases)
2. Das neue Release sollte sichtbar sein mit allen Binaries
3. Teste einen Download um sicherzustellen, dass die Binaries funktionieren
4. Überprüfe die Release Notes auf Vollständigkeit

## 🔧 Troubleshooting

### Release-Workflow schlägt fehl

1. **Überprüfe die Logs**: Gehe zu Actions und schaue dir die Fehlermeldungen an
2. **Häufige Probleme**:
   - Tests schlagen fehl → Fixe die Tests erst
   - Build-Fehler → Überprüfe Cross-Platform Kompatibilität
   - Tag bereits vorhanden → Siehe "Tag löschen" unten

### Tag löschen (falls nötig)

```bash
# Tag lokal löschen
git tag -d v1.2.0

# Tag remote löschen  
git push origin --delete v1.2.0

# GitHub Release manuell löschen:
# Gehe zu: https://github.com/mmuyakwa/WhatsApp-Console/releases
# Klicke auf das Release → Edit → Delete
```

### Hotfix Release

Für kritische Bugfixes:

```bash
# Hotfix direkt auf master
git checkout master
git pull origin master

# Bugfix committen
git add .
git commit -m "Hotfix: Fix critical authentication bug"
git push origin master

# Patch-Version erhöhen
git tag -a v1.2.1 -m "Hotfix v1.2.1: Fix critical authentication bug"
git push origin v1.2.1
```

## 📝 Best Practices

### Versioning

- **MAJOR** (v1.0.0 → v2.0.0): Breaking Changes, API-Änderungen
- **MINOR** (v1.0.0 → v1.1.0): Neue Features, backwards-kompatibel
- **PATCH** (v1.0.0 → v1.0.1): Bugfixes, Security-Patches

### Release Notes

Verwende aussagekräftige Commit-Messages:
```bash
git commit -m "Add: Support for group message broadcasting"
git commit -m "Fix: Memory leak in message handler" 
git commit -m "Update: Improve error messages for better UX"
```

### Pre-Releases

Für Beta/Alpha-Versionen:
```bash
# Beta-Release
git tag -a v1.2.0-beta.1 -m "Release v1.2.0-beta.1: Beta version with new features"
git push origin v1.2.0-beta.1
```

### Testing vor Release

```bash
# Lokale Tests
go test ./...

# Cross-Platform Build Test
GOOS=linux GOARCH=amd64 go build -o test-linux .
GOOS=windows GOARCH=amd64 go build -o test-windows.exe .
GOOS=darwin GOARCH=amd64 go build -o test-macos .

# Funktionalitäts-Test
./test-linux --version

# Aufräumen
rm test-*
```

## 📊 Nach dem Release

### Monitoring

1. **Download-Statistiken** auf GitHub überprüfen
2. **Issues/Bug Reports** überwachen
3. **User Feedback** in Discussions sammeln

### Wartung

1. **Dependabot PRs** zeitnah reviewen
2. **Security Alerts** ernst nehmen
3. **Performance-Monitoring** bei größeren Releases

---

## 🔗 Weiterführende Links

- [GitHub Releases](https://github.com/mmuyakwa/WhatsApp-Console/releases)
- [GitHub Actions](https://github.com/mmuyakwa/WhatsApp-Console/actions)  
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
