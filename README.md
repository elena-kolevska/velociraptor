# Webhook Rate Limiter Microservice

A lightweight interceptor service written in [Go](https://golang.org/) meant to serve as a rate limiter for webhooks from external services.  
Originally I built it for Twilio Chat webhooks, because the project I was working on didn't have another way to sort chats by last activity, and we didn't want to ping the API (a full blown Laravel app) for every  webhook.  

In its heart, the microservice uses Redis keys with a configurable expiration time. [TODO]

## Features:

### 1. Keeping track of conversation last activity time and date  

#### Objective:
Act as a rate-limiter to our [main API](https://docs.google.com/document/d/1tjzWGZtSdMOKT40Rj_2vmedq0WzFhk254rxKV7iV09U/edit#heading=h.og3cyed1lb2j) in the process of updating the last activity timestamp of conversations. 

#### Overview:
- Our Twilio chat system is [configured](https://www.twilio.com/console/chat/services/IS896ada4aceeb4058b230383eb3199be1/webhooks) in such a way that every time a new message is sent in a conversation, it sends a webhook to this microservice.  
(Read more details about [Twilio webhooks](https://www.twilio.com/docs/chat/webhook-events))
- Sending a webhook to the [main API](https://docs.google.com/document/d/1tjzWGZtSdMOKT40Rj_2vmedq0WzFhk254rxKV7iV09U/edit#heading=h.og3cyed1lb2j) for every single message exchanged on our plartform would strain the servers and main database, so we put this lightweight microservice in the middle to act as a rate limiter. 
- By making use of Redis it only allows updates to a channel [using the API](https://commune-x.postman.co/collections/94102-14df5d05-5487-4288-a96a-346b92d909bd?workspace=eacc830c-4864-4381-976c-d113c29881e3#54b46f5f-32a2-4699-b268-b92eae41ec9f) once per minute.


### 2. Recognise booking requests and create tasks in Asana
