// routes.go

package main

func initializeRoutes(cws *careWorkerServer) {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	cws.router.Use(cws.setUserStatus())

	// Handle the index route
	cws.router.GET("/q", cws.showIndexPage)

	// Group user related routes together
	cws.userRoutes = cws.router.Group("/u")
	{
		// Handle POST requests at /u/login
		cws.userRoutes.POST("/login", cws.ensureNotLoggedIn(), cws.performLogin)

		// Handle GET requests at /u/logout
		cws.userRoutes.GET("/logout", cws.ensureLoggedIn(), cws.logout)

		// Handle POST requests at /u/register
		cws.userRoutes.POST("/register", cws.ensureNotLoggedIn(), cws.register)

		// Handle POST requests at /u/register
		cws.userRoutes.POST("/register/salt", cws.ensureNotLoggedIn(), cws.registerSalt)
	}

	// Group article related routes together
	cws.articleRoutes = cws.router.Group("/article")
	{
		// Handle GET requests at /article/view/some_article_id
		cws.articleRoutes.GET("/view/:article_id", cws.getArticle)

		// Handle the GET requests at /article/create
		cws.articleRoutes.GET("/create", cws.ensureLoggedIn(), cws.showArticleCreationPage)

		// Handle POST requests at /article/create
		cws.articleRoutes.POST("/create", cws.ensureLoggedIn(), cws.createArticle)

		// Handle POST requests at /article/update
		cws.articleRoutes.POST("/update", cws.ensureLoggedIn(), cws.updateArticle)

		// Handle the GET requests at /article/delete/:id
		cws.articleRoutes.GET("/delete/:id", cws.ensureLoggedIn(), cws.deleteArticle)
	}
}
