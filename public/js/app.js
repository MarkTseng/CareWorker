var careworkerApp = angular.module('careworker', ['angularMoment', 'ngRoute', 'careworkerControllers', 'ngCookies', 'pascalprecht.translate'])

careworkerApp.config(['$routeProvider',
	function($routeProvider) {
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
			otherwise({
				redirectTo: '/questions'
			});
	}
]);

careworkerApp.run(function($rootScope, $location, AuthService, FlashService, SessionService) {
	var routesThatRequireAuth = ['/ask'];

	$rootScope.authenticated = SessionService.get('authenticated');
	$rootScope.username = SessionService.get('username');

	$rootScope.$on('$routeChangeStart', function (event, next, current) {
		FlashService.clear()
		if(_(routesThatRequireAuth).contains($location.path()) && !AuthService.isLoggedIn())
		{
			FlashService.show("Please login to continue");
			$location.path('/login');
		}
	});
});

careworkerApp.config(function($httpProvider) {
	var logsOutUserOn401 = function ($location, $q, SessionService, FlashService) {
		var success = function (res) {
			return res;
		}
		var error   = function (res) {
			if(res.status === 401) { // HTTP NotAuthorized
				SessionService.unset('authenticated')
				FlashService.show(res.data.Message);
				$location.path("/login");
				return $q.reject(res)
			} else {
				return $q.reject(res)
			}
		}

		return function(promise) {
			return promise.then(success, error)
		}
	}

	$httpProvider.responseInterceptors.push(logsOutUserOn401);
})

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
	function($rootScope, $http, $location, SessionService, FlashService) {

		var cacheSession = function (user) {
			SessionService.set('authenticated', true);
			SessionService.set('username', user.Username);
            $rootScope.authenticated = true;
			$rootScope.user = user;
		}

		var uncacheSession = function () {
			SessionService.unset('authenticated');
			SessionService.unset('username');
			$rootScope.authenticated = false;
			$rootScope.user = {};
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

			for ( var i = 0; i < arguments.length; i++) {
				s += arguments[i];
			}
			return SHA256.hex(s);

		}

		return {
			login: function (credentials, fn) {

				// Friendly vars
				var u = credentials.Username;
				var p = credentials.Password;

				credentials.Username = "";
				credentials.Password = "";

				// Get a salt for this session
				$http.post("/u/register/salt", {"username" : u})
					.success(function(user_salt) {
								// Produce the "Password" to send
								p = protect (u + p, user_salt.salt);
								//p = hash( p , user_salt.salt)

								// Try to login
								var login = $http.post("/u/login", {"username": u, "password": p, "salt": user_salt.salt});

								login.success(cacheSession);
								login.success(FlashService.clear);
								login.error(loginError);

								if ( typeof fn === "function" )
									login.success(fn);
					}
				)
			},
			logout: function (fn) {
				var logout =  $http.get("/u/logout");
				logout.success(uncacheSession);

				if ( typeof fn === "function" )
					logout.success(fn);

			},
			register: function (credentials, fn) {

				// Friendly vars
				var u = credentials.username;
				var p = credentials.password;

				credentials.username = "";
				credentials.password = "";

				var s = protect(Date.now(), Math.random());

				// Produce the "Password" to send
				p = protect(u + p, s);

				var register = $http.post("/u/register", {"username" : u, "password" : p, "salt" : s });

				if ( typeof fn === "function")
					register.success(fn);

			},
			isLoggedIn: function () {
				return SessionService.get('authenticated');
			},
			currentUser: function () {
				if ( this.isLoggedIn() )
				{
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
			if ( $rootScope.flashn-- == 0 )
				$rootScope.flash = ""
		}
	}
});

careworkerApp.factory('Questions', ['$http',
	function ($http) {

		var urlBase = '/q'
		var store	= []
		var f		= {};

		$http.defaults.headers.common['Accept'] = 'application/json';
		var p = function (data) {
			this.success = function (fn) {
				fn(data)
			}
		}

		f.List = function () {
			return $http.get(urlBase)
				.success(function (data) {
					store = data;
				});
		}

		f.Get = function (id, force) {
			if ( !force ) {
				for ( var i = 0; i < store.length; i++ )
				{
					if ( store[i].ID == id )
						return new p(store[i]); 
				}
			}
			
			return $http.get(urlBase + '/' + id);
		}

		f.Insert = function (item) {
			return $http.post(urlBase, item)
				.success(function(data) {
					store.push(data);
				});
		}

		f.Update = function (item) {
			return $http.put(urlBase + '/' + item.ID, item)
				.success(function (data) {
					for ( var i = 0; i < store.length; i++ )
					{
						if ( store[i].ID == id )
							return store[i] = data;
					}
				});
		}

		f.Delete = function (id) {
			return $http.delete(urlBase + '/' + id)
				.success(function (data) {
					for ( var i = 0; i < store.length; i++ )
					{
						if ( store[i].ID == id )
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
		 link: function(scope, element, attrs) {
			var evHandler = function(e) {
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
		function link ( scope, element, attributes ) {
			console.log("Generating question...", attributes.question);
			Questions.Get(attributes.question).success(function(data) {
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
		function link ( scope, element, attributes ) {
			console.log("Generating comment...", attributes.question);
			Questions.Get(attributes.question).success(function(data) {
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
	return function(input) {
		if(input === 0)
			return "at least enter 15 characters";
		else if(input < 15)
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

// Custom validator based on expressions.
// see: https://docs.angularjs.org/guide/forms
careworkerApp.directive('wjValidationError', function () {
  return {
    require: 'ngModel',
    link: function (scope, elm, attrs, ctl) {
      scope.$watch(attrs['wjValidationError'], function (errorMsg) {
        elm[0].setCustomValidity(errorMsg);
        ctl.$setValidity('wjValidationError', errorMsg ? false : true);
      });
    }
  };
});
