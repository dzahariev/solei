sap.ui.define([
	'./BaseController',
	'sap/ui/model/json/JSONModel',
	'sap/ui/Device',
	'com/zahariev/solei/model/formatter'
], function (BaseController, JSONModel, Device, formatter) {
	"use strict";
	return BaseController.extend("com.zahariev.solei.controller.Home", {
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
		},

		onExit: function() {
			Object.keys(this.events).forEach(eventName => {
				sap.ui.getCore().getEventBus().unsubscribe("solei", eventName, this.events[eventName], this);
			});
		},

		handleRouteChanged: function(channel, eventId, pageData) {
			if (pageData.selectedPageKey === "home"){
				// Ignore for now
			}
		},
	});
});