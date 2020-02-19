# WHSubscriber
The official WebhookShareSystem-subscriber. 

### What is a WebhookShareSystem
Imagine you want to run a shell script after your favourite gets a new release, or a new docker image on dockehub is available. For example, you want to update your dockercontainer everytime a new image gets released. With this you'll face following problems:
- You don't have access right to the repository so you can't add a webhook
- You would have to setup a lot of complicated stuff
- You can't parse the payload of the webhook easily in bash

This WeShareSystem allows you to share your webhooks with others, to allow them to subscribe these webhooks and run actions if something happens on your repo. In addition if your favourite repository supports WeShareSystem you can run your own actions (eg udate the docker container, etc...)
