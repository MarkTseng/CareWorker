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
		cws.userRoutes.GET("/logout", cws.ensureLoggedIn(), cws.logout)
		cws.userRoutes.POST("/register", cws.ensureNotLoggedIn(), cws.register)
		cws.userRoutes.POST("/register/salt", cws.ensureNotLoggedIn(), cws.registerSalt)
		cws.userRoutes.POST("/login", cws.ensureNotLoggedIn(), cws.performLogin)
	}

	// Group article related routes together
	cws.articleRoutes = cws.router.Group("/article")
	{
		cws.articleRoutes.GET("/view/:article_id", cws.getArticle)
		cws.articleRoutes.GET("/create", cws.ensureLoggedIn(), cws.showArticleCreationPage)
		cws.articleRoutes.GET("/delete/:id", cws.ensureLoggedIn(), cws.deleteArticle)
		cws.articleRoutes.POST("/create", cws.ensureLoggedIn(), cws.createArticle)
		cws.articleRoutes.POST("/update", cws.ensureLoggedIn(), cws.updateArticle)
	}
}
