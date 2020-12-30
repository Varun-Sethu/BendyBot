# Bendy Bot

A relatively simple discord bot designed to track people's messages and generate random sentences. 

### Usage
To use Bendy Bot make sure you have an active bot registered with Discord, this can be done through their website. Attain the authentication token for your both and execute the following command in the home directory:
```bash
echo 'token' >> storage/data/authcode.txt
```

Upon the creation of the authentication token the bot can then be run through the following commands:

```bash
docker build -t mrme/bendybot .
docker run -d mrme/bendybot
```

### Commands
Bendy bot consists of the following commands:
- `yo bendy track @user` - Begins tracking a user
- `yo bendy endtrack @user` - Stops tracking a user
- `yo bendy generate @user` - Generates a sentence from the user's dictionary
- `yo bendy set-channel` - Changes the tracking channel to wherever the command was called from

