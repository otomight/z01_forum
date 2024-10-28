# Forum 

## Specifications 

### Project main goals

#### Communication between users 

In order for users to communicate between each other, they'll have to be able to create Posts and Comments which'll only be allowed to registered users.

I think displaying latest Posts first when accessing the home page is the more coherent course of actionas well as accessing the latest comments for each post first.

##### Posts

- Create Posts : Registered user should be able to write and publish a Post. 

- Filters : Registered users can filter the displayer post by

    - Categories (all users).
    - Own created Posts.
    - Liked Posts.
    - (Favorite Posts).

- Likes : Liking or disliking Posts is only available to registered users.

    - I'd prefer the user not to be able to like/dislik his own Posts.

- History : the user should be able to access the history of 

    - His published Posts.
    - His liked/disliked Posts.

- (Modify/Delete his own Posts).

##### Comments

- Comment existing Posts is for registered users only.

- Like/dislike comments.

    - Same here, I'd prefer the user not to be able to like/dislike his own comment.

- (Modify/Delete his own comments).

- History : user should be able to see all his published comments.

#### Authentication

Clients must be able to register as new users on the Forum by inputting their credentials as well as log in.

The use of cookies is essential to allow each user to have only one opened session :
- Each session must contain an expiration date.

##### Register

- User registration :

    - Username.
    - Email.
    - Password.
    - (FirstName).
    - (LastName).
    - (BirthDate).

The web application must be able to check if the Email address is present in the database and if all credentials are correct.

Users should be able to register using at least Facebook, Github and Google.

##### Log in

- User Login :

    - Username/Email.
    - Password.

When logging, the Forum must be able to check if the entered password is the same as the one provided in the database.

Users should also be able to Log in using at least Facebook, Github and Goolge.

#### Moderation

Moderation is essential for a Forum to control the displayed content so that it respects the Term and Condition of Use (TCUs).

The content must be appropriate : respectful, relevant; not obscene, illegal or insulting.

To check the conformity of the Forum content we'll need a moderation system with several levels.

##### Guest

- Guests : unregistered users that can only read the Posts and Comments, they also can search for a Post topic using the search bar or by Category.

##### Basic User 

- Users : Able to create, comment, like or dislike posts:

    - They 'll also be able to dislike or modify/delete their own Posts/Comments.

    - They should be able to request to be promoted to moderator to the admin.

##### Moderator

- Moderators : they have the same options than the basic users but with some special features :

    - Report or delete users/Posts or comments to admin.
    - Report/Ban users.
    - Request a user moderator promotion to the admin.

##### Administrator

- Admin : Above the moderator, he manages the technical details required for running the Forum :

    - Promote/demote user to moderator or back to basic user.
    - Receive reports from moderators.
    - Delete Posts/Comments/users 
    - Manage categories (add/delete them).

#### Security

##### HTTPS protocol

##### Rate Limiting

##### Password encryption

##### Database encryption

##### Unique client session cookies


## Epics/user stories/Back log : global features of the project 

![Epics/user stories](image.png)

- Epic : Create a web application which allow users to communicate via posts/comments.

- User stories :
    - Be able to register/Login.
    - Be able to make a post.
    - Be able to choose one/multiple categories fo the post .
    - Be able to like/dislike a post/comment .
    - Be able to see the history of created/liked-disliked posts and comments.
    - Be able to see all posts and comments from selected categories.

- Back log :
    - Authentication : email/username/password + error handling
    - Post creation : registered users only :
        - Choose one/multiple categories for the post. 
    - Like-dislike posts/comments : registered users only.
    - SQLite : 
        - History of created posts/comments and likes-dislikes : registered users only.
        - History of all posts/comment from selected categories : ALL users 
    
    ### User stories 

The project started with the determination of the Main Epic and the different user stories.
The Epic represents the major objective of the project and the user stories the different features needed to respond to the customer's demands.
These user stories are prioritized to determine which ones we should do first : 
- High priority
- Moderate priority
- Low priority 

They also contain validation criteria to ensure that the functionality is correctly implemented.

![Product Owner user](image-1.png)

![Product owner user types](image-2.png)

## Data dictionary 

It details the database we use in our informative system. It’s used by devs & database admins, and is really helpful in its précise data description.

![User/Email verification data](image-3.png)

![Session data](image-5.png)

![Post data](image-4.png)

## MERISE Method

This method is useful to link the different elements of our project together. 
To follow this method we used the free Looping software which is intuitive and simple to apprehend.

###  Conceptual Data Model

Entities are represented by the Yellow rectangles and contain their different attributes.
Blue links associate the entities to each other when needed, they sometime have to contain some attributes as well, they also indicate how many times 2 Entities can be associated to each other : 
 - One to one (1,1) : Each entity occurrence is only linked to one occurrence of another entity.
- Zero or One to Many  (0,n) : One entity occurrence can be (or not) associated with several occurrences of another entity.


![CMD](image-6.png)

### Logical Data Structure 

It’s used to detail entity relations. The links defined in the CDM are translated into primary and foreign keys to establish the relations between tables : 
- Primary keys : Underlined & bold, unique attribute identifier.
- Foreign keys : Underline & bold but blue, reference to another table primary key.

![LSD](image-7.png)

### Physical Data Model

It’s the last step of the data modelisation in which the LDS is translated to data structures specific to the used database.

![PDM](image-8.png)

![PDM 2](image-9.png)

## Wireframe  

The Wireframe is like the Blueprint of the project, it can be divided in subgroups depending on its accuracy : 
- Low-fidelity : the most basic one
    - Layout
    - Navigation
    - Informative architecture 
- Mid-fidelity : 
    - Mapping out core functionalities/ key interactions
    - Adding annotation/content
- High-fidelity : It’s like a early mockup of the project 
    - Interactive/visual design elements
    - Fonts/colors/logos
    
We agree on the fact  the first draft isn’t definitive, it’ll evolve as we find more features to add and as we modify the structure.

![WireFrame timeline](image-22.png)

### Low-Level

#### Mobile

![Mobile home page](image-10.png)

![Mobile sign up/log in](image-11.png)

![Mobile New Post](image-14.png)

#### Desktop

![Desktop home page](image-12.png)

![Desktop signup/log in](image-13.png)

![Desktop New Post](image-15.png)


### Medium-Level

#### Mobile

![Mobile user parameters](image-16.png)

![Mobile Posts/Likes History](image-17.png)

![Mobile saved to Favs](image-18.png)

#### Desktop

![Desktop user parameters](image-19.png)

![Desktop Posts/Likes History](image-20.png)

![Desktop saved to Favs](image-21.png)


### High-Level 

![Functionnal Wireframe](image-24.png)


## Technologies used 

### Back-End : Go

- The language we used fot he back-end part of the project is GoLang (developped by google) for several reasons :
    - It was imposed 
    - Simple, efficient & secure
        - Performance : Compiles directly in binary while being easy to write
        - Concurrency : Excellent handling of concurrency with goroutines (essential for simultaneous interactions between forum's users).

#### Server 

- We can only use standard Go packages for the project, the easiest way to build a server adapted to a forum would have been to do a Websocket server for real-time communication between clients and server. It'd have been the best option beacause we need to be able to get instant updates about the last Posts, comments and likes notifications.
- We had to think about another way to mimick the websocket :
    - Long polling 
        - The client repeatedly makes requests to the server which responds immediatly with and update or waits for a certain period before responding.
        - Server has to repeatedly establish new client-server connections for each client (latency and more resources consumption)
    - SSE (Server-Sent Events)
        - Maintains a single HTTP connection open for each client : allows the server to push updates to the client whenever needed
        - Only sends data when there's an update 

### Front-End : JS, HTML, CSS

- JS : Interactive User Interfaces creation
    - DOM (Document Object Model) manipulation
    - API calls (fetch)
    - Dynamic interactivity 
    - Handle SSE connections :
        - Create an EventSource : initiate EventSource Object with ther server endpoint URL that'll be provinding SSE data 
        - Listening for Events : The client listens for events sent by the server and each event triggers a JS function.
        - Updating Page : JS updates the webpage with new data by modifying the DOM to show the new messages, updates, notifications 
- Html : Standard web taging structure
- CSS :  Responsive Layout and interfaces design 

### Docker 

Docker is used to package and run an application in an isolated environment called a container, developers can work in a standardized environment. 
Containers are lightweight and contain everything needed to run the application.
It's possible to share the containers while working : ensures that everyone gets an application that works in the same way.

What we can do with Docker :
- Develop our application and its supporting components using containers
- The container becoles the unit for distributing and testing our application
- When ready : allows to deploy our application into our production environment as a container and an orchestrated service

#### Docker architecture 

![Docker architecure](image-23.png)

Docker uses a Client-Server architecture :
- The Client : talks to the Docker Deamon
- The Deamon builds, runs and distributes our containers (can also communicate with other deamons)

The Docker Client and Deamon can run on the same system or we can connect a Client to a remote Deamon. They communicate using REST API(over Unix sockets or network interface)


### Database : SQLite 

- For the Database, we were asked to use SQLite, a C-language library that implements small, fast, self-contained high-reliability, full-featured SQL database engine.SQLite is built in every smartphone and most computers.
- SQLite file format is stable, cross-platform. SQLite database files are commonly used as containers to transfer rich content between systems & has a long term archival format for data 
- SQLite is not comparable to client-server SQL database engines because it provides local data storage for individual applications and devices. It competes with the Fopen() function.
- It works well as database engine for low to medium traffic websites (100K to 500K request/day).