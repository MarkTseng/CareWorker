var careworkerControllers = angular.module('careworkerControllers', ['ngMaterial', 'ngMessages', 'ngCookies']);

careworkerControllers.controller('QuestionListCtrl', ['$scope', 'Questions',
	function ($scope, Questions) {
		Questions.List()
			.then(function successCallback(data) {
				$scope.questions = data;
                $scope.time = new Date();
				console.log($scope.time);
			}, function errorCallback(error) {
				$scope.error = 'Unable to load questions';
			});
	}
]);

careworkerControllers.controller('RegisterCtrl', ['$scope', '$http', '$location', 'AuthService', 'FlashService',
	function ($scope, $http, $location, AuthService, FlashService) {
		var credentials = {
			confirmPassword: "",
			email: "",
			username: "",
			password: "",
		};
		$scope.credentials = credentials;

		$scope.processForm = function () {
			$scope.errorComment = ""
			// Default to a non error state
			var err = false;

			// Make sure they have entered a password that matches
			if ($scope.credentials.password != $scope.credentials.confirmPassword) {
				err = true;
				return;
			}
			if (err)
				return;

			// debug
			//console.log($scope.credentials)
			//return;

			// We don't need this to be passed along
			delete $scope.credentials.confirmPassword;

			// Register the user and redirect them to the login page
			AuthService.register($scope.credentials, function successCallback(response) {
				$location.path("/login");
			}, function errorCallback(response) {
				FlashService.show(response.data[0].Message);
			});

			// Make sure we wipe out the credentials
			$scope.credentials = credentials;

		}
	}
]);

careworkerControllers.controller('LoginCtrl', ['$scope', '$http', '$location', 'AuthService',
	function ($scope, $http, $location, AuthService) {
		var credentials = { "Password": "", "Salt": "", "Email": "" }

		$scope.credentials = credentials
		$scope.processForm = function () {

			// Log the user in and direct them tot he home page
			AuthService.login($scope.credentials, function () {
				$location.path("/");
			});

			// Make sure we wipe out the credentials
			$scope.credentials = credentials
		}
	}
]);

careworkerControllers.controller('ResetPasswordCtrl', ['$scope', '$http', '$routeParams', '$location', 'AuthService', 'FlashService',
	function ($scope, $http, $routeParams, $location, AuthService, FlashService) {
		var credentials = { "Password": "", "email": "", "resetCode":"" }
        $scope.credentials = credentials

        credentials["email"] = $routeParams.email
        credentials["resetCode"] = $routeParams.resetcode
		$scope.processForm = function () {
            //console.log("reset password")
			// We don't need this to be passed along
			delete $scope.credentials.confirmPassword;

			AuthService.resetPassword($scope.credentials, function successCallback(response) {
				$location.path("/login");
			}, function errorCallback(response) {
				FlashService.show(response.data[0].Message);
			});
		}
	}
]);

careworkerControllers.controller('ForgotPasswordCtrl', ['$scope', '$http', '$location', 'AuthService', 'FlashService',
	function ($scope, $http, $location, AuthService, FlashService) {
		var credentials = { "email": "" }
        $scope.credentials = credentials
        FlashService.clear()

		$scope.processForm = function () {
            $http.defaults.headers.common['Accept'] = 'application/json';
			$scope.errorComment = ""
			// Default to a non error state
			var err = false;

			// Make sure they have entered a password that matches
			if ($scope.credentials.email === "") {
				err = true;
				return;
			}

			if (err)
				return;

	        $http.get('/u/forgotpassword/' + $scope.credentials.email).then(function successCallback(response) {
				$location.path("/login");
		    }, function errorCallback(response) {
				FlashService.show(response.data[0].Message);
			});
		}
	}
]);


careworkerControllers.controller('QuestionAskCtrl', ['$scope', '$http', '$window', '$sce', '$location',
	function ($scope, $http, $window, $sce, $location) {
		var question = { "markdown": "", "title": "", "tags": "" }

		$scope.question = question;

		$scope.md2Html = function () {
			var src = $scope.question.markdown || ""
			var html = $window.marked(src);
			$scope.htmlSafe = $sce.trustAsHtml(html);
		}

		$scope.processForm = function () {

			// Remove any previous error statements
			$scope.error = {}


			// Default to a non error state
			var err = false;

			if ($scope.question.markdown.length < 10) {
				$scope.error.markdown = "Your question must be 10 characters or more."
				err = true;
			}

			if ($scope.question.title.length == 0) {
				$scope.error.title = "You must enter a title."
				err = true;
			}

			if ($scope.question.tags.length == 0) {
				$scope.error.tags = "You must have at least one tag."
				err = true;
			}

			if (err) {
				return;
			}

			$http.defaults.headers.common['Accept'] = 'application/json';
			$http({
				method: 'POST',
				url: '/article/create',
				data: { title: $scope.question.title, body: $scope.question.markdown, Tags: $scope.question.tags.split(' ') }
			}).then(function (response) {
				// TODO: this should be a JSON response
				$location.path("/questions/" + response.data);
			});
			// TODO: Failure
		}

	}
]);

careworkerControllers.controller('carrerCtrl', ['$scope', '$http', '$window', '$sce', '$location',
	function ($scope, $http, $window, $sce, $location) {
		var carrer = { "serviceType": "保母服務", "location": "", "salary": "", "body": "" }

		// service types
		$scope.serviceTypes = [
			"看護服務",
			"清潔服務",
		];
		$scope.serviceType = $scope.serviceTypes[0];

		// request types
		$scope.requestServices = ["到府托育","全日托","日托","半日托","臨托","其他"];
		$scope.requestServicesSelected = [];

		$scope.requestServicesToggle = function (item, list) {
			var idx = list.indexOf(item);
			if (idx > -1) {
				list.splice(idx, 1);
			}
			else {
				list.push(item);
			}
		};

		$scope.exists = function (item, list) {
			return list.indexOf(item) > -1;
		};

		$scope.carrer = carrer;
		$scope.carrer['requestService'] = $scope.requestServicesSelected

		// process HR form table
		$scope.processForm = function () {

			// Remove any previous error statements
			$scope.error = {}


			// Default to a non error state
			var err = false;

			if ($scope.job.title.length == 0) {
				$scope.error.title = "You must enter a title."
				err = true;
			}

			if ($scope.job.location.length == 0) {
				$scope.error.tags = "You must enter a location."
				err = true;
			}

			if ($scope.job.salary.length == 0) {
				$scope.error.tags = "You must enter a salary."
				err = true;
			}

			if ($scope.job.body.length == 0) {
				$scope.error.tags = "You must enter a content."
				err = true;
			}

			if (err) {
				return;
			}

			$http.defaults.headers.common['Accept'] = 'application/json';
			$http({
				method: 'POST',
				url: '/article/create',
				data: { title: $scope.job.title, location: $scope.job.location, salary: $scope.job.salary, body: $scope.job.body }
			}).then(function (response) {
				// TODO: this should be a JSON response
				$location.path("/questions/" + response.data);
			});
			// TODO: Failure
		}

	}
]);


careworkerControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http', '$window', '$sce',
	function ($scope, $routeParams, $http, $window, $sce) {
		$scope.comment = { "Body": "" }
		$scope.response = { "Body": "" }
		$http.defaults.headers.common['Accept'] = 'application/json';
		$http.get('/article/view/' + $routeParams.questionId).then(function (response) {
			$scope.question = response.data[0];
			//console.log(response.data)
		});

		$scope.voteUp = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question._id + '/vote/up',
				data: {}
			}).then(function (response) {
				$scope.question.Upvotes = response.data.Upvotes
			});
		}

		$scope.voteDown = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question._id + '/vote/down',
				data: {}
			}).then(function (response) {
				$scope.question.Downvotes = response.data.Downvotes
			});
		}

		$scope.markdown = "";
		$scope.md2Html = function () {
			var src = $scope.response.Body || ""
			$scope.html = $window.marked(src);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		// Can a comment have this own controller and it's own scope?
		$scope.processComment = function () {
			delete $scope.errorComment;

			var err = false;

			if ($scope.comment.Body.length < 15) {
				$scope.errorComment = "Your comment must be at least 15 characters"
				err = true;
			}

			if (err) return;

			$http({
				method: 'post',
				url: '/q/' + $scope.question._id + '/comment/',
				data: $scope.comment
			}).then(function (data) {
				delete $scope.scomment;
				$scope.question.Comments.push(data);
			});
		}

		$scope.processForm = function () {
			//console.log($scope.response.Body);
			delete $scope.errorMarkdown;

			var err = false;

			if ($scope.response.Body.length < 5) {
				$scope.errorMarkdown = "Your response must be 5 characters or more."
				err = true;
			}


			if (err) {
				return;
			}

			$http({
				method: 'post',
				url: '/q/' + $scope.question._id + '/response/',
				data: $scope.response
			}).then(function (data) {
				$scope.question.Responses.push(data);
			});
		}
	}
]);

careworkerControllers.controller('carrerDetailCtrl', ['$scope', '$routeParams', '$http', '$window', '$sce',
	function ($scope, $routeParams, $http, $window, $sce) {
		$http.defaults.headers.common['Accept'] = 'application/json';
		$http.get('/article/view/' + $routeParams.jobId).then(function (response) {
			$scope.job = response.data[0];
		});
	}
]);

careworkerControllers.controller('ProfileCtrl', ['$scope', '$routeParams', '$http', '$window', '$sce', '$location', 'AuthService', 'SessionService', 'FlashService',
	function ($scope, $routeParams, $http, $window, $sce, $location, AuthService, SessionService, FlashService) {
		var profile = {
			birthday: "",
			userId: "",
			idtype: "",
			birthdayUTC: new Date(2010,0,1),
			phone: "",
			city: "",
			district: "",
			name: "",
			gender: "",
			street: "",
			zipcode: "",
			license: "",
			jobbrew: ""
        };
        
        $http.defaults.headers.common['Accept'] = 'application/json';
        var user = JSON.parse(SessionService.get('user'));
		$http.get('/u/profile/' + user.UserId).then(function (response) {
			profile = response.data[0];
			SessionService.set('profile', response.data[0]);
    
            profile['name'] = user.Username
            $scope.updateBirthday = function () {
                var DateBirthday = new Date($scope.profile.birthdayUTC);
                // UTC to GMT(taiwan) + 8 
                DateBirthday.setHours(DateBirthday.getHours() + 8);
                $scope.profile.birthday = DateBirthday.toISOString().slice(0,10);
            }
            $scope.updateZipcode = function () {
                if ($scope.profile.city && $scope.profile.district) {
                    $scope.profile.zipcode = $scope.areas[$scope.profile.city][$scope.profile.district];
                } else {
                    $scope.profile.zipcode = undefined;
                }
            };

            $scope.updateDistrict = function () {
                $scope.districts = Object.keys($scope.areas[$scope.profile.city])
            }

            $scope.areas = {
                基隆市: {
                    仁愛區: "200", 信義區: "201", 中正區: "202", 中山區: "203", 安樂區: "204", 暖暖區: "205",
                    七堵區: "206"
                },
                臺北市: {
                    中正區: "100", 大同區: "103", 中山區: "104", 松山區: "105", 大安區: "106", 萬華區: "108",
                    信義區: "110", 士林區: "111", 北投區: "112", 內湖區: "114", 南港區: "115", 文山區: "116"
                },
                新北市: {
                    萬里區: "207", 金山區: "208", 板橋區: "220", 汐止區: "221", 深坑區: "222", 石碇區: "223",
                    瑞芳區: "224", 平溪區: "226", 雙溪區: "227", 貢寮區: "228", 新店區: "231", 坪林區: "232",
                    烏來區: "233", 永和區: "234", 中和區: "235", 土城區: "236", 三峽區: "237", 樹林區: "238",
                    鶯歌區: "239", 三重區: "241", 新莊區: "242", 泰山區: "243", 林口區: "244", 蘆洲區: "247",
                    五股區: "248", 八里區: "249", 淡水區: "251", 三芝區: "252", 石門區: "253"
                },
                宜蘭縣: {
                    宜蘭市: "260", 頭城鎮: "261", 礁溪鄉: "262", 壯圍鄉: "263", 員山鄉: "264", 羅東鎮: "265",
                    三星鄉: "266", 大同鄉: "267", 五結鄉: "268", 冬山鄉: "269", 蘇澳鎮: "270", 南澳鄉: "272",
                    釣魚臺列嶼: "290"
                },
                新竹市: { 東區: "300", 北區: "300", 香山區: "300" },
                新竹縣: {
                    竹北市: "302", 湖口鄉: "303", 新豐鄉: "304", 新埔鎮: "305", 關西鎮: "306", 芎林鄉: "307",
                    寶山鄉: "308", 竹東鎮: "310", 五峰鄉: "311", 橫山鄉: "312", 尖石鄉: "313", 北埔鄉: "314",
                    峨眉鄉: "315"
                },
                桃園市: {
                    中壢區: "320", 平鎮區: "324", 龍潭區: "325", 楊梅區: "326", 新屋區: "327", 觀音區: "328",
                    桃園區: "330", 龜山區: "333", 八德區: "334", 大溪區: "335", 復興區: "336", 大園區: "337",
                    蘆竹區: "338"
                },
                苗栗縣: {
                    竹南鎮: "350", 頭份市: "351", 三灣鄉: "352", 南庄鄉: "353", 獅潭鄉: "354", 後龍鎮: "356",
                    通霄鎮: "357", 苑裡鎮: "358", 苗栗市: "360", 造橋鄉: "361", 頭屋鄉: "362", 公館鄉: "363",
                    大湖鄉: "364", 泰安鄉: "365", 銅鑼鄉: "366", 三義鄉: "367", 西湖鄉: "368", 卓蘭鎮: "369"
                },
                臺中市: {
                    中區: "400", 東區: "401", 南區: "402", 西區: "403", 北區: "404", 北屯區: "406", 西屯區: "407",
                    南屯區: "408", 太平區: "411", 大里區: "412", 霧峰區: "413", 烏日區: "414", 豐原區: "420",
                    后里區: "421", 石岡區: "422", 東勢區: "423", 和平區: "424", 新社區: "426", 潭子區: "427",
                    大雅區: "428", 神岡區: "429", 大肚區: "432", 沙鹿區: "433", 龍井區: "434", 梧棲區: "435",
                    清水區: "436", 大甲區: "437", 外埔區: "438", 大安區: "439"
                },
                彰化縣: {
                    彰化市: "500", 芬園鄉: "502", 花壇鄉: "503", 秀水鄉: "504", 鹿港鎮: "505", 福興鄉: "506",
                    線西鄉: "507", 和美鎮: "508", 伸港鄉: "509", 員林市: "510", 社頭鄉: "511", 永靖鄉: "512",
                    埔心鄉: "513", 溪湖鎮: "514", 大村鄉: "515", 埔鹽鄉: "516", 田中鎮: "520", 北斗鎮: "521",
                    田尾鄉: "522", 埤頭鄉: "523", 溪州鄉: "524", 竹塘鄉: "525", 二林鎮: "526", 大城鄉: "527",
                    芳苑鄉: "528", 二水鄉: "530"
                },
                南投縣: {
                    南投市: "540", 中寮鄉: "541", 草屯鎮: "542", 國姓鄉: "544", 埔里鎮: "545", 仁愛鄉: "546",
                    名間鄉: "551", 集集鎮: "552", 水里鄉: "553", 魚池鄉: "555", 信義鄉: "556", 竹山鎮: "557",
                    鹿谷鄉: "558"
                },
                嘉義市: { 東區: "600", 西區: "600" },
                嘉義縣: {
                    番路鄉: "602", 梅山鄉: "603", 竹崎鄉: "604", 阿里山: "605", 中埔鄉: "606", 大埔鄉: "607",
                    水上鄉: "608", 鹿草鄉: "611", 太保市: "612", 朴子市: "613", 東石鄉: "614", 六腳鄉: "615",
                    新港鄉: "616", 民雄鄉: "621", 大林鎮: "622", 溪口鄉: "623", 義竹鄉: "624", 布袋鎮: "625"
                },
                雲林縣: {
                    斗南鎮: "630", 大埤鄉: "631", 虎尾鎮: "632", 土庫鎮: "633", 褒忠鄉: "634", 東勢鄉: "635",
                    臺西鄉: "636", 崙背鄉: "637", 麥寮鄉: "638", 斗六市: "640", 林內鄉: "643", 古坑鄉: "646",
                    莿桐鄉: "647", 西螺鎮: "648", 二崙鄉: "649", 北港鎮: "651", 水林鄉: "652", 口湖鄉: "653",
                    四湖鄉: "654", 元長鄉: "655"
                },
                臺南市: {
                    中西區: "700", 東區: "701", 南區: "702", 北區: "704", 安平區: "708", 安南區: "709", 永康區: "710",
                    歸仁區: "711", 新化區: "712", 左鎮區: "713", 玉井區: "714", 楠西區: "715", 南化區: "716",
                    仁德區: "717", 關廟區: "718", 龍崎區: "719", 官田區: "720", 麻豆區: "721", 佳里區: "722",
                    西港區: "723", 七股區: "724", 將軍區: "725", 學甲區: "726", 北門區: "727", 新營區: "730",
                    後壁區: "731", 白河區: "732", 東山區: "733", 六甲區: "734", 下營區: "735", 柳營區: "736",
                    鹽水區: "737", 善化區: "741", 大內區: "742", 山上區: "743", 新市區: "744", 安定區: "745"
                },
                高雄市: {
                    新興區: "800", 前金區: "801", 苓雅區: "802", 鹽埕區: "803", 鼓山區: "804", 旗津區: "805",
                    前鎮區: "806", 三民區: "807", 楠梓區: "811", 小港區: "812", 左營區: "813", 仁武區: "814", 大社區: "815",
                    東沙群島: "817", 南沙群島: "819", 岡山區: "820", 路竹區: "821", 阿蓮區: "822", 田寮區: "823", 燕巢區: "824",
                    橋頭區: "825", 梓官區: "826", 彌陀區: "827", 永安區: "828", 湖內區: "829", 鳳山區: "830", 大寮區: "831",
                    林園區: "832", 鳥松區: "833", 大樹區: "840", 旗山區: "842", 美濃區: "843", 六龜區: "844", 內門區: "845",
                    杉林區: "846", 甲仙區: "847", 桃源區: "848", 那瑪夏區: "849", 茂林區: "851", 茄萣區: "852"
                },
                屏東縣: {
                    屏東市: "900", 三地門鄉: "901", 霧臺鄉: "902", 瑪家鄉: "903", 九如鄉: "904", 里港鄉: "905",
                    高樹鄉: "906", 鹽埔鄉: "907", 長治鄉: "908", 麟洛鄉: "909", 竹田鄉: "911", 內埔鄉: "912",
                    萬丹鄉: "913", 潮州鎮: "920", 泰武鄉: "921", 來義鄉: "922", 萬巒鄉: "923", 崁頂鄉: "924",
                    新埤鄉: "925", 南州鄉: "926", 林邊鄉: "927", 東港鎮: "928", 琉球鄉: "929", 佳冬鄉: "931",
                    新園鄉: "932", 枋寮鄉: "940", 枋山鄉: "941", 春日鄉: "942", 獅子鄉: "943", 車城鄉: "944",
                    牡丹鄉: "945", 恆春鎮: "946", 滿州鄉: "947"
                },
                臺東縣: {
                    臺東市: "950", 綠島鄉: "951", 蘭嶼鄉: "952", 延平鄉: "953", 卑南鄉: "954", 鹿野鄉: "955",
                    關山鎮: "956", 海端鄉: "957", 池上鄉: "958", 東河鄉: "959", 成功鎮: "961", 長濱鄉: "962",
                    太麻里鄉: "963", 金峰鄉: "964", 大武鄉: "965", 達仁鄉: "966"
                },
                花蓮縣: {
                    花蓮市: "970", 新城鄉: "971", 秀林鄉: "972", 吉安鄉: "973", 壽豐鄉: "974", 鳳林鎮: "975",
                    光復鄉: "976", 豐濱鄉: "977", 瑞穗鄉: "978", 萬榮鄉: "979", 玉里鎮: "981", 卓溪鄉: "982",
                    富里鄉: "983"
                },
                金門縣: {
                    金沙鎮: "890", 金湖鎮: "891", 金寧鄉: "892", 金城鎮: "893", 烈嶼鄉: "894", 烏坵鄉: "896"
                },
                連江縣: { 南竿鄉: "209", 北竿鄉: "210", 莒光鄉: "211", 東引鄉: "212" },
                澎湖縣: {
                    馬公市: "880", 西嶼鄉: "881", 望安鄉: "882", 七美鄉: "883", 白沙鄉: "884", 湖西鄉: "885"
                }
            };

            $scope.citys = Object.keys($scope.areas)
            $scope.districts = Object.keys($scope.areas[$scope.citys[0]])
            $scope.profile = profile;   
            // fixed onload error
            if(profile.city != "")
                $scope.districts = Object.keys($scope.areas[profile.city])


            if($scope.profile.birthday != "")
            {
                var strBirthday = new String($scope.profile.birthday);
                $scope.profile.birthdayUTC = new Date(strBirthday.split("-"));
            }
            $scope.processForm = function () {
                if($scope.profile.birthday === "")
                {
                    $scope.updateBirthday();
                }

                AuthService.profile($scope.profile, function successCallback(response) {
                    $location.path("/profile");
                }, function errorCallback(response) {
                    $scope.errorComment = response.data[0].Message
					FlashService.show(response.data[0].Message);
                });
            }
		});
    }
]);
