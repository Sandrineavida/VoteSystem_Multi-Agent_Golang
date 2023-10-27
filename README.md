
# Types

Les types importants sont définis dans le fichier  `types\types.go`.

## Une récap pour tous les types 

| Type          | Description                                                                                                  |
|---------------|--------------------------------------------------------------------------------------------------------------|
| Alternative   | Un type pour représenter les alternatives (candidats).                                                        |
| Profile       | Un type pour représenter les profils (liste des préférences des électeurs).                                   |
| Count         | Un type pour représenter les résultats (map des scores(value) des candidats(key)).                           |
| Ballot        | Structure pour représenter un bulletin de vote.                                                              |
| VoteRequest   | Structure pour représenter une demande de vote.                                                              |
| NewBallotRequest | Structure pour représenter une demande pour créer un nouveau bulletin de vote.                          |
| ResultVoteRequest | Structure pour représenter une demande de résultat de vote.                                              |
| BaseResponse  | Structure de type de base pour les réponses, contenant un indicateur de succès et un message éventuel.        |
| NewBallotResponse | Structure pour représenter la réponse à la requête /new_ballot.                                         |
| VoteResponse  | Structure pour représenter la réponse à la requête /vote.                                                     |
| ResultVoteResponse | Structure pour représenter la réponse à la requête /result.                                             |

## Détails de  `Ballot`

| Nom de param   | Type de param    | Description de param                                                                         |
|----------------|------------------|---------------------------------------------------------------------------------------------|
| ID             | string           | ID unique d'un ballot                                                                       |
| Rule           | string           | Règle de vote                                                                               |
| Deadline       | time.Time        | Date limite de vote                                                                         |
| Votes          | [][]Alternative  | Tous les votes (Profile)                                                                    |
| HasCounted     | bool             | Indique si le dépouillement a été effectué.                                                 |
| VoterIDs       | []string         | Liste des ID des votants ayant déjà voté                                                    |
| Alts           | int              | Nombre de candidats                                                                         |
| TieBreak       | []Alternative    | L'ordre prédéfini des candidats utilisé quand on a un "tie"                                 |
| Result         | Alternative      | Résultat du vote (Un gagnant)                                                               |
| Thresholds     | []int            | Dans le cadre d'Approval, le nombre de votes pour chaque électeur                           |
| Seuil          | int              | Seuil dans le cadre d'Approval (qui doit <= Alts int)                                       |

## Détails de `VoteRequest`

| Nom de param | Type de param | Description de param                   | Nom en json |
|--------------|---------------|----------------------------------------|-------------|
| AgentID      | string        | ID de l'électeur                       | agent-id    |
| BallotID     | string        | ID du ballot                           | ballot-id   |
| Prefs        | []Alternative | Préférences de l'électeur               | prefs       |
| Options      | int           | Nombre d'options considerées            | options     |

## Détails de `NewBallotRequest`

| Nom de param | Type de param | Description de param                                                       | Nom en json      |
|--------------|---------------|----------------------------------------------------------------------------|------------------|
| Rule         | string        | Règle de vote                                                              | rule             |
| Deadline     | time.Time     | Date limite pour voter                                                     | deadline         |
| VoterIDs     | []string      | Liste des IDs des votants                                                  | voter-ids        |
| Alts         | int           | Nombre de candidats                                                        | #alts            |
| TieBreak     | []Alternative | L'ordre prédéfini des candidats utilisé quand on a un "tie"                 | tie-break        |
| Seuil        | int           | Seuil dans le cadre d'Approval (qui doit <= Alts int)                       | seuil            |

## Détails de `ResultVoteRequest`

| Nom de param | Type de param | Description de param                           | Nom en json  |
|--------------|---------------|------------------------------------------------|--------------|
| BallotID     | string        | ID du ballot                                   | ballot-id    |


## Détails de `BaseResponse`

| Nom de param | Type de param | Description de param                                                    | Nom en json |
|--------------|---------------|-------------------------------------------------------------------------|-------------|
| Success      | bool          | Indique si la demande a réussi                                          | success     |
| Message      | string        | Message d'erreur ou de succès potentiel                                 | message     |

## Détails de `NewBallotResponse`

| Nom de param | Type de param | Description de param                                                     | Nom en json       |
|--------------|---------------|--------------------------------------------------------------------------|-------------------|
| -            | BaseResponse | Intégration de la structure de réponse de base                           | -                 |
| BallotID     | string        | ID du bulletin nouvellement créé                                         | ballot_id         |

## Détails de `VoteResponse`

| Nom de param | Type de param | Description de param                                                     | Nom en json |
|--------------|---------------|--------------------------------------------------------------------------|-------------|
| -            | BaseResponse  | Intégration de la structure de réponse de base                           | -           |

## Détails de `ResultVoteResponse`

| Nom de param | Type de param | Description de param                                                     | Nom en json |
|--------------|---------------|--------------------------------------------------------------------------|-------------|
| -            | BaseResponse  | Intégration de la structure de réponse de base                           | -           |
| Winner       | Alternative   | Le gagnant du vote                                                       | winner      |


# Serveur

Le fichier `server.go` situé dans `vote\restserveragent` implémente les fonctionnalités côté serveur.

## Structure d'un agent serveur `RestServerAgent`

| Champ    | Type                    | Description             |
|----------|-------------------------|-------------------------|
| -        | sync.Mutex              | Mutex pour synchronization|
| id       | string                  | ID d'un agent de serveur|
| reqCount | int                     | Compteur de requêtes    |
| ballots  | map[string]types.Ballot | map de ballots          |
| addr     | string                  | Adresse                 |

## Méthodes définies dans `server.go`

| Méthode                         | Utilisation                                                     |
| ------------------------------- | --------------------------------------------------------------- |
| NewRestServerAgent              | Instancier un agent de serveur (type RestServerAgent)           |
| checkMethod                     | Test de la méthode                                              |
| decodeNewBallotRequest          | Decoder les requetes de type NewBallotRequest                   |
| decodeVoteRequest               | Decoder les requetes de type VoteRequest                        |
| decodeResultVoteRequest         | Decoder les requetes de type ResultVoteRequest                  |
| doReqNewBallot                  | Cote serveur: traiter la requete pour creer un nouveau ballot   |
| doVote                          | Cote serceur: traiter la requete pour voter                     |
| doReqcount                      | Donner le nombre de requêtes reçues par le serveur              |
| doReqresult                     | Cote serveur: traiter la requete pour obtenir le résultat du vote |
| calculateResultCheck            | Essayer de Calculer le résultat du vote (toutes les 10 secondes) |
| Start                           | Lancer le serveur                                               |

## Endpoint et Handler Function

| Endpoint    | Handler Function |
|-------------|------------------|
| /vote       | doVote           |
| /count      | doReqcount       |
| /result     | doReqresult      |
| /new_ballot | doReqNewBallot   |

# Client

Le fichier `client.go` situé dans `vote\restclientagent` implémente les fonctionnalités côté client. 

*N.B.*: Nous avons *2* types de clients :
- Des clients qui demandent au serveur de créer un ballot
- Des clients qui agissent en tant qu'électeurs

## Structure d'un agent client `RestClientAgent`

| Champ      | Type                   | Description                                                                 |
|------------|------------------------|-----------------------------------------------------------------------------|
| id         | string                 | ID de l'agent                                                               |
| url        | string                 | URL du serveur                                                              |
| prefs      | []types.Alternative    | Préférences de l'agent (Si c'est un électeur)                               |
| seuil      | int                    | Valeur uniquement pour Approval (Si c'est un électeur)                      |
| options    | int                    | Nombre d'options considerées (Si c'est un électeur)                         |
| nCandidat  | int                    | Nombre de candidats                                                         |

## Méthodes définies pour RestClientAgent

| Méthode                      | Utilisation                                                                                                           |
|------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| NewRestClientAgent           | Instancier un agent client                                                                                            |
| treatVoteResponse            | Traiter la réponse du serveur pour la requête de voter                                                                |
| treatNewBallotResponse       | Traiter la réponse du serveur pour la requête de créer un nouveau ballot                                              |
| treatResultVoteResponse      | Traiter la réponse du serveur pour la requête de résultat                                                             |
| treatCountResponse           | Traiter la réponse du serveur pour la requête de compter le nombre de requêtes                                        |
| doNewBallotRequest           | Demander au serveur de créer un nouveau ballot                                                                        |
| doVoteRequest                | Demander de voter sur un ballot donné                                                                                 |
| doResultVoteRequest          | Demander le résultat d'un ballot donné                                                                                |
| doCountRequest               | Demander le nombre total de requêtes envoyées au serveur                                                              |
| CreateBallotStart            | Fonction de démarrage pour l'agent qui génère les ballots                                                             |
| Start                        | Fonction de démarrage pour l'agent qui est électeur, contient l'action de voter, de demander les résultats et de demander le nombre d'appels au serveur |
