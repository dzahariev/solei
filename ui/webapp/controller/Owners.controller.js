sap.ui.define([
	'./BaseController',
	'sap/ui/model/json/JSONModel',
	'sap/ui/Device',
	'com/zahariev/solei/model/formatter'
], function (BaseController, JSONModel, Device, formatter) {
	"use strict";
	return BaseController.extend("com.zahariev.solei.controller.Owners", {
		formatter: formatter,

		onInit: function () {
			// Subscribe on events
			this.events = {
				RouteChanged: this.handleRouteChanged,
			};

			Object.keys(this.events).forEach(eventName => {
				sap.ui.getCore().getEventBus().subscribe("solei", eventName, this.events[eventName], this);
			});

			// Define device model
			var oViewModel = new JSONModel({
				isPhone : Device.system.phone
			});
			this.setModel(oViewModel, "view");
			Device.media.attachHandler(function (oDevice) {
				this.getModel("view").setProperty("/isPhone", oDevice.name === "Phone");
			}.bind(this));

			// Define owner model
			this.ownerModel = new JSONModel({
				owners: null,
			})
			this.ownerModel.setSizeLimit(Number.MAX_VALUE);
			this.setModel(this.ownerModel, "owner");
			this.loadData()
			this.getView().setModel(this.ownerModel, "owner");
		},

		onExit: function() {
			Object.keys(this.events).forEach(eventName => {
				sap.ui.getCore().getEventBus().unsubscribe("solei", eventName, this.events[eventName], this);
			});
		},

		handleRouteChanged: function(channel, eventId, pageData) {
			if (pageData.selectedPageKey === "home"){
				this.loadData()
			}
		},

		loadData: async function () {
			const ownerModel = this.getModel("owner")
			const token = await this.getOwnerComponent().getToken()
			var strResponse = "";
			jQuery.ajax({
				url: '/api/owner',
				type: "GET",
				beforeSend: function (xhr) {
					xhr.setRequestHeader('Authorization', `Bearer ${token}`);
				},
				success: function(response) {
					strResponse = response;
				},
				async:false
			  });
	
			ownerModel.setProperty("/owners", strResponse.data);
			this.getView().getModel("owner").refresh(true);
		}
	});
});