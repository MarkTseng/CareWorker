// routes.go

package main

func initializeRoutes(cws *careWorkerServer) {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	cws.router.Use(cws.setUserStatus())

	// Handle the index route
	//cws.router.GET("/", cws.showIndexPage)

	// Handle css static file
	cws.router.Static("/css", "public/static/css")

	// Handle css static file
	cws.router.Static("/photos", "public/static/photos")

	// Group user related routes together
	cws.userRoutes = cws.router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		cws.userRoutes.GET("/login", cws.ensureNotLoggedIn(), cws.showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		cws.userRoutes.POST("/login", cws.ensureNotLoggedIn(), cws.performLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		cws.userRoutes.GET("/logout", cws.ensureLoggedIn(), cws.logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		cws.userRoutes.GET("/register", cws.ensureNotLoggedIn(), cws.showRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		cws.userRoutes.POST("/register", cws.ensureNotLoggedIn(), cws.register)
	}

	// Group article related routes together
	cws.articleRoutes = cws.router.Group("/article")
	{
		// Handle GET requests at /article/view/some_article_id
		cws.articleRoutes.GET("/view/:article_id", cws.getArticle)

		// Handle the GET requests at /article/create
		// Show the article creation page
		// Ensure that the user is logged in by using the middleware
		cws.articleRoutes.GET("/create", cws.ensureLoggedIn(), cws.showArticleCreationPage)

		// Handle POST requests at /article/create
		// Ensure that the user is logged in by using the middleware
		cws.articleRoutes.POST("/create", cws.ensureLoggedIn(), cws.createArticle)

		// Handle POST requests at /article/update
		// Ensure that the user is logged in by using the middleware
		cws.articleRoutes.POST("/update", cws.ensureLoggedIn(), cws.updateArticle)

		// Handle the GET requests at /article/delete/:id
		// delete article
		// Ensure that the user is logged in by using the middleware
		cws.articleRoutes.GET("/delete/:id", cws.ensureLoggedIn(), cws.deleteArticle)
	}
}
