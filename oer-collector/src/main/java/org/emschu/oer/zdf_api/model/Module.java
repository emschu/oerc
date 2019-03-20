
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

import java.util.ArrayList;
import java.util.List;
import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

public class Module {

    @SerializedName("shorttext-headline")
    @Expose
    private String shorttextHeadline;
    @SerializedName("shorttext-text")
    @Expose
    private String shorttextText;
    @SerializedName("profile")
    @Expose
    private String profile;
    @SerializedName("teaser")
    @Expose
    private List<Teaser> teaser = new ArrayList<Teaser>();
    @SerializedName("title")
    @Expose
    private String title;
    @SerializedName("layout")
    @Expose
    private String layout;
    @SerializedName("maxBroadcasts")
    @Expose
    private Integer maxBroadcasts;
    @SerializedName("showBroadcastRepeats")
    @Expose
    private Boolean showBroadcastRepeats;
    @SerializedName("homeTvService")
    @Expose
    private String homeTvService;
    @SerializedName("showTitle")
    @Expose
    private Boolean showTitle;

    public String getShorttextHeadline() {
        return shorttextHeadline;
    }

    public void setShorttextHeadline(String shorttextHeadline) {
        this.shorttextHeadline = shorttextHeadline;
    }

    public String getShorttextText() {
        return shorttextText;
    }

    public void setShorttextText(String shorttextText) {
        this.shorttextText = shorttextText;
    }

    public String getProfile() {
        return profile;
    }

    public void setProfile(String profile) {
        this.profile = profile;
    }

    public List<Teaser> getTeaser() {
        return teaser;
    }

    public void setTeaser(List<Teaser> teaser) {
        this.teaser = teaser;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getLayout() {
        return layout;
    }

    public void setLayout(String layout) {
        this.layout = layout;
    }

    public Integer getMaxBroadcasts() {
        return maxBroadcasts;
    }

    public void setMaxBroadcasts(Integer maxBroadcasts) {
        this.maxBroadcasts = maxBroadcasts;
    }

    public Boolean getShowBroadcastRepeats() {
        return showBroadcastRepeats;
    }

    public void setShowBroadcastRepeats(Boolean showBroadcastRepeats) {
        this.showBroadcastRepeats = showBroadcastRepeats;
    }

    public String getHomeTvService() {
        return homeTvService;
    }

    public void setHomeTvService(String homeTvService) {
        this.homeTvService = homeTvService;
    }

    public Boolean getShowTitle() {
        return showTitle;
    }

    public void setShowTitle(Boolean showTitle) {
        this.showTitle = showTitle;
    }

}
