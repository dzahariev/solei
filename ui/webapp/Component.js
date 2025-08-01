sap.ui.define([
	"sap/ui/core/library",
	"sap/ui/core/UIComponent",
	"./model/models",
	"sap/ui/model/json/JSONModel",
	"sap/ui/core/routing/History",
	"sap/ui/Device",
	"sap/ui/model/resource/ResourceModel"
], function(library, UIComponent, models, JSONModel, History, Device) {
	"use strict";

	return UIComponent.extend("com.zahariev.solei.Component", {
		metadata: {
			manifest: "json",
			interfaces: [library.IAsyncContentCreation]
		},

		init: async function () {
			// call the init function of the parent
			UIComponent.prototype.init.apply(this, arguments);

			// define auth model and init keycloak
			this.authModel = new JSONModel({
				keycloak: null,
				username: "Unknown",
			})
			var Keycloak;
			if (window.Keycloak) {
				Keycloak = window.Keycloak;
			}

			this.authModel.setSizeLimit(Number.MAX_VALUE);
			const keycloak = new Keycloak('/keycloak-cfg/keycloak.json');
			const authenticated = await keycloak.init({ onLoad: 'login-required', checkLoginIframe: false })
			if(authenticated){
				this.authModel.setProperty("/keycloak", keycloak);
				this.authModel.setProperty("/username", keycloak.tokenParsed?.preferred_username);
			} else {
				keycloak.logout();
			}
			this.setModel(this.authModel, "authModel");
			
			// set the device model
			this.setModel(models.createDeviceModel(), "device");

			// set the menu model
			const sideModel = this.getModel('sideContent')
			this.setModel(sideModel, "side");

			// create the views based on the url/hash
			this.getRouter().initialize();
		},

		getToken: async function () {
			const authModel = this.getModel("authModel");
			const keycloak = authModel.getProperty("/keycloak");
			if(keycloak.isTokenExpired()){
				try {
					const tokenUpdated = await keycloak.updateToken(1800)
					if(tokenUpdated) {
						authModel.setProperty("/username", keycloak.tokenParsed?.preferred_username);
					} else {
						keycloak.logout();
					}
				} catch (error) {
					keycloak.logout();
				}
			}
			return keycloak.token;
		},

		hasRole: function (role) {
			const authModel = this.getModel("authModel");
			const keycloak = authModel.getProperty("/keycloak");
			return keycloak.hasRealmRole(role)
		},

		getContentDensityClass: function () {
			if (!this._sContentDensityClass) {
				if (!Device.support.touch){
					this._sContentDensityClass = "sapUiSizeCompact";
				} else {
					this._sContentDensityClass = "sapUiSizeCozy";
				}
			}
			return this._sContentDensityClass;
		}
	});
});