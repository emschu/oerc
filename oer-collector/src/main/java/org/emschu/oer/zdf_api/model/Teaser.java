
package org.emschu.oer.zdf_api.model;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

public class Teaser {

    @SerializedName("http://zdf.de/rels/target")
    @Expose
    private HttpZdfDeRelsTarget httpZdfDeRelsTarget;
    @SerializedName("headerTitle")
    @Expose
    private String headerTitle;
    @SerializedName("headerDecorationText")
    @Expose
    private String headerDecorationText;
    @SerializedName("displayBrandLogo")
    @Expose
    private Boolean displayBrandLogo;
    @SerializedName("brandLogoColor")
    @Expose
    private String brandLogoColor;
    @SerializedName("profile")
    @Expose
    private String profile;

    public HttpZdfDeRelsTarget getHttpZdfDeRelsTarget() {
        return httpZdfDeRelsTarget;
    }

    public void setHttpZdfDeRelsTarget(HttpZdfDeRelsTarget httpZdfDeRelsTarget) {
        this.httpZdfDeRelsTarget = httpZdfDeRelsTarget;
    }

    public String getHeaderTitle() {
        return headerTitle;
    }

    public void setHeaderTitle(String headerTitle) {
        this.headerTitle = headerTitle;
    }

    public String getHeaderDecorationText() {
        return headerDecorationText;
    }

    public void setHeaderDecorationText(String headerDecorationText) {
        this.headerDecorationText = headerDecorationText;
    }

    public Boolean getDisplayBrandLogo() {
        return displayBrandLogo;
    }

    public void setDisplayBrandLogo(Boolean displayBrandLogo) {
        this.displayBrandLogo = displayBrandLogo;
    }

    public String getBrandLogoColor() {
        return brandLogoColor;
    }

    public void setBrandLogoColor(String brandLogoColor) {
        this.brandLogoColor = brandLogoColor;
    }

    public String getProfile() {
        return profile;
    }

    public void setProfile(String profile) {
        this.profile = profile;
    }

}
