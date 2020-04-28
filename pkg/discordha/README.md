# DiscordHA

DiscordHA (Discord High Available) is library on to be used together with discord-go to deploy Discord bots in high availability.
It relies on Etcd as a locking system to prevent events of being received twice, this works in a first locked principle for now.

Discord HA is not meant for sharding but enables Discord bots to have multiple replicas.