### What is this

Metro is a simulation of a complex transport system. At its current state it only cares about train transport.

### Things missing

- There should be a brain that observes and interferes with the system that may be eventually controlled by a math/ai model. Besides it should be the main governor of the data and analytics.
- People. There should be people interacting with the system. Boarding trains, having destinations they need to go, etc. They should have names and some sentiment on how they feel based on time waiting for trains, hwo confortable they are (depends on the capacity of the trains and how much is filled)
- Real data. In previous instances of this project, I got real data of the train system of Santo Domingo which created all the stations in their relative position. Eventually I would like to be able to have at least 2 to 3 cities that the system can simulate.
- A model that can summarize the given actions of the system and make it into a newspaper somehow.
- There should be a score of how the simulation is going based on how the people are feeling about the system.

- The simulation should have the concept of time, meaning that trains need to be in X station by X time. Tenjin should work towards those times being as truth as possible. No later not that earlier arrivals, of course that part is for when tenjin will actually use a model.
- Everything related to the UI/UX of the application, but that will be addressed later on.

### What the final deliverable looks like

This simulation is supposed to be interacted with by a touch screen in the real life.
The possible interactions are:

- User can click on a train and check the trains current data (where its heading, the amount of passengers inside, current velocity)
- User can click on a station and check the station data (people waiting inline for train, next train incoming, how many people has served)
- User can check the newspaper to figure out a summary of that day's simulation, including the people's sentiments for example.
