<div align="center">
<article style="display: flex; flex-direction: column; align-items: center; justify-content: center;">
  <h1 style="width: 100%; text-align: center;">Product store</h1>
  <p>Product store. Task description is <a href="task.md">here</a></p>
</article>
</div>

# 🔥 Run
Download:
```sh
git clone https://github.com/ernur-eskermes/product-store.git 
mv .env.example .env
make run
```
Stop:
```sh
make down
```
Running multiple instances with a load balancer:
```sh
make deploy
```

# What is going on there
What does the client do:
- sending address to download products (for example, run project <a href="github.com/ernur-eskermes/lead-csv-service/">lead-csv-service</a>)
- requesting products with pagination functionality (page, size)
- requesting products continuously, simulate endless loading (based on bidirectional-stream)

What does the server do:
- going to the url, download and save the csv-file, extract the products from it, save to your database
- giving products with pagination functionality (page, size)
- giving away products in stream

Services are raised in docker-compose, the client starts and executes requests with output to the console, and then exits.  
The server continues to work  