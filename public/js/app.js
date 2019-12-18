var careworkerApp = angular.module('careworker', ['careworkerControllers', 'ngMaterial', 'ngMessages', 'ngRoute', 'ngCookies', 'pascalprecht.translate']);

careworkerApp.config(['$routeProvider',
	function ($routeProvider) {
		$routeProvider.
			when('/questions', {
				templateUrl: 'partials/question-list.html',
				controller: 'QuestionListCtrl'
			}).
			when('/questions/:questionId', {
				templateUrl: 'partials/question-detail.html',
				controller: 'QuestionDetailCtrl'
			}).
			when('/job/:jobId', {
				templateUrl: 'partials/job-detail.html',
				controller: 'JobDetailCtrl'
			}).
			when('/ask', {
				templateUrl: 'partials/question-post.html',
				controller: 'QuestionAskCtrl'
			}).
			when('/job', {
				templateUrl: 'partials/job-post.html',
				controller: 'JobCtrl'
			}).
			when('/register', {
				templateUrl: 'partials/register.html',
				controller: 'RegisterCtrl'
			}).
			when('/login', {
				templateUrl: 'partials/login.html',
				controller: 'LoginCtrl'
			}).
			when('/profile', {
				templateUrl: 'partials/profile.html',
				controller: 'ProfileCtrl'
			}).
			otherwise({
				redirectTo: '/questions'
			});
	}
]);

careworkerApp.run(function ($rootScope, $location, AuthService, FlashService, SessionService) {
	var routesThatRequireAuth = ['/ask'];

	$rootScope.authenticated = SessionService.get('authenticated');
	$rootScope.username = SessionService.get('username');
	$rootScope.userId = SessionService.get('userId');

	$rootScope.$on('$routeChangeStart', function (event, next, current) {
		FlashService.clear()
		if (_(routesThatRequireAuth).contains($location.path()) && !AuthService.isLoggedIn()) {
			FlashService.show("Please login to continue");
			$location.path('/login');
		}
	});
});

careworkerApp.config(['$httpProvider', function ($httpProvider) {
	var logsOutUserOn401 = function ($location, $q, SessionService, FlashService) {
		var success = function (res) {
			return res;
		}
		var error = function (res) {
			if (res.status === 401) { // HTTP NotAuthorized
				SessionService.unset('authenticated')
				FlashService.show(res.data.Message);
				$location.path("/login");
				return $q.reject(res)
			} else {
				return $q.reject(res)
			}
		}

		return function (promise) {
			return promise.then(success, error)
		}
	}

	$httpProvider.interceptors.push(logsOutUserOn401);
}])

careworkerApp.factory("SessionService", function () {
	return {
		get: function (key) {
			return sessionStorage.getItem(key);
		},
		set: function (key, val) {
			return sessionStorage.setItem(key, val);
		},
		unset: function (key) {
			return sessionStorage.removeItem(key);
		}
	}
});

careworkerApp.factory("AuthService", ['$rootScope', '$http', '$location', 'SessionService', 'FlashService',
	function ($rootScope, $http, $location, SessionService, FlashService) {

		var cacheSession = function (response) {
			SessionService.set('authenticated', true);
			SessionService.set('username', response.data.Username);
			SessionService.set('userId', response.data.UserId);
			$rootScope.authenticated = true;
			$rootScope.user = response.data;
			$rootScope.username = response.data.Username;
			$rootScope.userId = response.data.UserId;
		}

		var uncacheSession = function () {
			SessionService.unset('authenticated');
			SessionService.unset('username');
			SessionService.unset('userId');
			$rootScope.authenticated = false;
			$rootScope.user = {};
			$rootScope.username = "";
			$rootScope.userId = "";
		}

		var loginError = function (res) {
			FlashService.show(res.Message);
		}

		var protect = function (secret, salt) {
			var SHA256 = new Hashes.SHA256;
			return SHA256.hex(salt + SHA256.hex(secret));
		}

		var hash = function () {
			var s = ""
			var SHA256 = new Hashes.SHA256;

			for (var i = 0; i < arguments.length; i++) {
				s += arguments[i];
			}
			return SHA256.hex(s);

		}

		return {
			login: function (credentials, fn) {

				// Friendly vars
				var e = credentials.Email;
				var p = credentials.Password;

				credentials.Email = "";
				credentials.Password = "";

				// Get a salt for this session
				$http.post("/u/register/salt", { "email": e })
					.then(function (response) {
						// Produce the "Password" to send
						p = protect(e + p, response.data.salt);
						//p = hash( p , response.data.salt)

						// Try to login
						var login = $http.post("/u/login", { "email": e, "password": p, "salt": response.data.salt });

						login.then(cacheSession);
						login.then(FlashService.clear);
						login.then(null, loginError);

						if (typeof fn === "function")
							login.then(fn);
					}
					)
			},
			logout: function (fn) {
				var logout = $http.get("/u/logout");
				logout.then(uncacheSession);

				if (typeof fn === "function")
					logout.then(fn);

			},
			register: function (credentials, successfn, errorfn) {

				// Friendly vars
				var e = credentials.email;
				var p = credentials.password;

				credentials.email = "";
				credentials.password = "";

				var s = protect(Date.now(), Math.random());

				// Produce the "Password" to send
				p = protect(e + p, s);

				//var register = $http.post("/u/register", {"email" : e, "password" : p, "salt" : s });
				var register = $http.post("/u/register", {
					"email": e,
					"password": p,
					"salt": s,
					"username": credentials.username,
				});
				if ((typeof successfn === "function") && (typeof errorfn === "function"))
					register.then(successfn, errorfn);

			},
			profile: function (profile, successfn, errorfn) {
                
		        profile['userId'] = SessionService.get('userId');
                console.log(profile)
				var profile = $http.post("/u/profile", profile);
				if ((typeof successfn === "function") && (typeof errorfn === "function"))
					profile.then(successfn, errorfn);

			},
			isLoggedIn: function () {
                // to do check session login
				return SessionService.get('authenticated');
			},
			currentUser: function () {
				if (this.isLoggedIn()) {
					return SessionService.get('user');
				}

				return {};
			}
		}
	}
]);

careworkerApp.factory("FlashService", function ($rootScope) {
	return {
		show: function (msg) {
			$rootScope.flashn = 1;
			$rootScope.flash = msg
		},
		clear: function () {
			if ($rootScope.flashn-- == 0)
				$rootScope.flash = ""
		}
	}
});

careworkerApp.factory('Questions', ['$http',
	function ($http) {

		var urlBase = '/q'
		var store = []
		var f = {};

		$http.defaults.headers.common['Accept'] = 'application/json';
		var p = function (data) {
			this.success = function (fn) {
				fn(data)
			}
		}

		f.List = function () {
			return $http.get(urlBase)
				.then(function (response) {
					store = response.data;
					// return value for then, next chain.
					return response.data;
				});
		}

		f.Get = function (id, force) {
			if (!force) {
				for (var i = 0; i < store.length; i++) {
					if (store[i].ID == id)
						return new p(store[i]);
				}
			}

			return $http.get(urlBase + '/' + id);
		}

		f.Insert = function (item) {
			return $http.post(urlBase, item)
				.then(function (response) {
					store.push(response.data);
				});
		}

		f.Update = function (item) {
			return $http.put(urlBase + '/' + item.ID, item)
				.then(function (response) {
					for (var i = 0; i < store.length; i++) {
						if (store[i].ID == id)
							return store[i] = response.data;
					}
				});
		}

		f.Delete = function (id) {
			return $http.delete(urlBase + '/' + id)
				.then(function (response) {
					for (var i = 0; i < store.length; i++) {
						if (store[i].ID == id)
							return store.splice(i, 1);
					}
				})
		}

		return f;
	}
]);

careworkerApp.directive('careworkerLogout', function (AuthService) {
	return {
		restrict: 'A',
		link: function (scope, element, attrs) {
			var evHandler = function (e) {
				e.preventDefault;
				AuthService.logout();
				return false;
			}

			element.on ? element.on('click', evHandler) : element.bind('click', evHandler);
		}
	}
});

careworkerApp.directive('question', ['Questions',
	function (Questions) {
		function link(scope, element, attributes) {
			console.log("Generating question...", attributes.question);
			Questions.Get(attributes.question).then(function (data) {
				console.log(data)
			})
		}

		return {
			restruct: 'A',
			link: link
		}
	}
]);

careworkerApp.directive('comment', ['Questions',
	function (Questions) {
		function link(scope, element, attributes) {
			console.log("Generating comment...", attributes.question);
			Questions.Get(attributes.question).then(function (data) {
				console.log(data)
			})
		}

		return {
			restruct: 'A',
			link: link
		}
	}
]);

careworkerApp.filter('commentremark', function () {
	return function (input) {
		if (input === 0)
			return "at least enter 15 characters";
		else if (input < 15)
			return "" + 15 - input + " more to go..."
		else
			return 600 - input + " characters left"

	}
});

careworkerApp.config(function ($translateProvider) {
	// configures staticFilesLoader
	$translateProvider.useStaticFilesLoader({
		prefix: 'data/locale-',
		suffix: '.json'
	});
	$translateProvider.preferredLanguage('zh');
});

careworkerApp.controller('Ctrl', ['$translate', '$scope', function ($translate, $scope) {
	$scope.changeLanguage = function (langKey) {
		$translate.use(langKey);
	};
}]);

careworkerApp.directive('pwCheck', function () {
	return {
		require: 'ngModel',
		scope: {
			password: '=match'
		},
		link: function (scope, element, attrs, ngModel) {
			scope.$watch('password', function () {
				ngModel.$setValidity('matchError', element.val() === scope.password);
			})
			element.on('keyup', function () {
				scope.$apply(function () {
					ngModel.$setValidity('matchError', element.val() === scope.password);
				})
			})
		}
	}
});

