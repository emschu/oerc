
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

import java.util.List;
import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

public class ZdfBroadcast {

    @SerializedName("playoutId")
    @Expose
    private String playoutId;
    @SerializedName("fsk")
    @Expose
    private Object fsk;
    @SerializedName("airtimeBegin")
    @Expose
    private String airtimeBegin;
    @SerializedName("airtimeEnd")
    @Expose
    private String airtimeEnd;
    @SerializedName("airtimeDate")
    @Expose
    private String airtimeDate;
    @SerializedName("effectiveAirtimeBegin")
    @Expose
    private String effectiveAirtimeBegin;
    @SerializedName("effectiveAirtimeEnd")
    @Expose
    private String effectiveAirtimeEnd;
    @SerializedName("audioComments")
    @Expose
    private Boolean audioComments;
    @SerializedName("blackwhite")
    @Expose
    private Boolean blackwhite;
    @SerializedName("zdfTvService")
    @Expose
    private String tvService;
    @SerializedName("tvServiceId")
    @Expose
    private String tvServiceId;
    @SerializedName("caption")
    @Expose
    private Boolean caption;
    @SerializedName("dolbyDigital51")
    @Expose
    private Boolean dolbyDigital51;
    @SerializedName("dolbySurround")
    @Expose
    private Boolean dolbySurround;
    @SerializedName("dualChannel")
    @Expose
    private Boolean dualChannel;
    @SerializedName("duration")
    @Expose
    private Long duration;
    @SerializedName("partDuration")
    @Expose
    private Long partDuration;
    @SerializedName("foreignLangWithCaption")
    @Expose
    private Boolean foreignLangWithCaption;
    @SerializedName("hd")
    @Expose
    private Boolean hd;
    @SerializedName("posId")
    @Expose
    private String posId;
    @SerializedName("partId")
    @Expose
    private Long partId;
    @SerializedName("live")
    @Expose
    private Boolean live;
    @SerializedName("livestream")
    @Expose
    private Boolean livestream;
    @SerializedName("modified")
    @Expose
    private String modified;
    @SerializedName("mono")
    @Expose
    private Boolean mono;
    @SerializedName("newAirtime")
    @Expose
    private Boolean newAirtime;
    @SerializedName("newProgramData")
    @Expose
    private Boolean newProgramData;
    @SerializedName("pharosId")
    @Expose
    private String pharosId;
    @SerializedName("signLanguage")
    @Expose
    private Boolean signLanguage;
    @SerializedName("stereo")
    @Expose
    private Boolean stereo;
    @SerializedName("visibleFrom")
    @Expose
    private String visibleFrom;
    @SerializedName("visibleTo")
    @Expose
    private String visibleTo;
    @SerializedName("vpsBegin")
    @Expose
    private String vpsBegin;
    @SerializedName("widescreen16_9")
    @Expose
    private Boolean widescreen169;
    @SerializedName("withChat")
    @Expose
    private Boolean withChat;
    @SerializedName("http://zdf.de/rels/cmdm/broadcasts-parts")
    @Expose
    private Object httpZdfDeRelsCmdmBroadcastsParts;
    @SerializedName("profile")
    @Expose
    private String profile;
    @SerializedName("self")
    @Expose
    private String self;
    @SerializedName("http://zdf.de/rels/tvservice")
    @Expose
    private ZdfTvService zdfTvService;
    @SerializedName("http://zdf.de/rels/cmdm/programme-item")
    @Expose
    private String httpZdfDeRelsCmdmProgrammeItem;
    @SerializedName("subtitle")
    @Expose
    private String subtitle;
    @SerializedName("text")
    @Expose
    private String text;
    @SerializedName("title")
    @Expose
    private String title;
    @SerializedName("subheadline")
    @Expose
    private String subheadline;
    @SerializedName("primaryBrandId")
    @Expose
    private String primaryBrandId;
    @SerializedName("primaryBrand")
    @Expose
    private String primaryBrand;
    @SerializedName("brandIds")
    @Expose
    private List<String> brandIds = null;
    @SerializedName("brandNames")
    @Expose
    private List<String> brandNames = null;
    @SerializedName("http://zdf.de/rels/image")
    @Expose
    private ZdfImage zdfImage;
    @SerializedName("onlineFrom")
    @Expose
    private Object onlineFrom;
    @SerializedName("onlineTo")
    @Expose
    private Object onlineTo;
    @SerializedName("geolocationVOD")
    @Expose
    private Object geolocationVOD;
    @SerializedName("geolocationLivestream")
    @Expose
    private Object geolocationLivestream;
    @SerializedName("youtubeRight")
    @Expose
    private Boolean youtubeRight;
    @SerializedName("onlineFirst")
    @Expose
    private Boolean onlineFirst;
    @SerializedName("http://zdf.de/rels/cmdm/series")
    @Expose
    private Object httpZdfDeRelsCmdmSeries;
    @SerializedName("http://zdf.de/rels/cmdm/season")
    @Expose
    private Object httpZdfDeRelsCmdmSeason;

    public String getPlayoutId() {
        return playoutId;
    }

    public void setPlayoutId(String playoutId) {
        this.playoutId = playoutId;
    }

    public ZdfBroadcast withPlayoutId(String playoutId) {
        this.playoutId = playoutId;
        return this;
    }

    public Object getFsk() {
        return fsk;
    }

    public void setFsk(Object fsk) {
        this.fsk = fsk;
    }

    public ZdfBroadcast withFsk(Object fsk) {
        this.fsk = fsk;
        return this;
    }

    public String getAirtimeBegin() {
        return airtimeBegin;
    }

    public void setAirtimeBegin(String airtimeBegin) {
        this.airtimeBegin = airtimeBegin;
    }

    public ZdfBroadcast withAirtimeBegin(String airtimeBegin) {
        this.airtimeBegin = airtimeBegin;
        return this;
    }

    public String getAirtimeEnd() {
        return airtimeEnd;
    }

    public void setAirtimeEnd(String airtimeEnd) {
        this.airtimeEnd = airtimeEnd;
    }

    public ZdfBroadcast withAirtimeEnd(String airtimeEnd) {
        this.airtimeEnd = airtimeEnd;
        return this;
    }

    public String getAirtimeDate() {
        return airtimeDate;
    }

    public void setAirtimeDate(String airtimeDate) {
        this.airtimeDate = airtimeDate;
    }

    public ZdfBroadcast withAirtimeDate(String airtimeDate) {
        this.airtimeDate = airtimeDate;
        return this;
    }

    public String getEffectiveAirtimeBegin() {
        return effectiveAirtimeBegin;
    }

    public void setEffectiveAirtimeBegin(String effectiveAirtimeBegin) {
        this.effectiveAirtimeBegin = effectiveAirtimeBegin;
    }

    public ZdfBroadcast withEffectiveAirtimeBegin(String effectiveAirtimeBegin) {
        this.effectiveAirtimeBegin = effectiveAirtimeBegin;
        return this;
    }

    public String getEffectiveAirtimeEnd() {
        return effectiveAirtimeEnd;
    }

    public void setEffectiveAirtimeEnd(String effectiveAirtimeEnd) {
        this.effectiveAirtimeEnd = effectiveAirtimeEnd;
    }

    public ZdfBroadcast withEffectiveAirtimeEnd(String effectiveAirtimeEnd) {
        this.effectiveAirtimeEnd = effectiveAirtimeEnd;
        return this;
    }

    public Boolean getAudioComments() {
        return audioComments;
    }

    public void setAudioComments(Boolean audioComments) {
        this.audioComments = audioComments;
    }

    public ZdfBroadcast withAudioComments(Boolean audioComments) {
        this.audioComments = audioComments;
        return this;
    }

    public Boolean getBlackwhite() {
        return blackwhite;
    }

    public void setBlackwhite(Boolean blackwhite) {
        this.blackwhite = blackwhite;
    }

    public ZdfBroadcast withBlackwhite(Boolean blackwhite) {
        this.blackwhite = blackwhite;
        return this;
    }

    public String getTvServiceId() {
        return tvServiceId;
    }

    public void setTvServiceId(String tvServiceId) {
        this.tvServiceId = tvServiceId;
    }

    public ZdfBroadcast withTvServiceId(String tvServiceId) {
        this.tvServiceId = tvServiceId;
        return this;
    }

    public Boolean getCaption() {
        return caption;
    }

    public void setCaption(Boolean caption) {
        this.caption = caption;
    }

    public ZdfBroadcast withCaption(Boolean caption) {
        this.caption = caption;
        return this;
    }

    public Boolean getDolbyDigital51() {
        return dolbyDigital51;
    }

    public void setDolbyDigital51(Boolean dolbyDigital51) {
        this.dolbyDigital51 = dolbyDigital51;
    }

    public ZdfBroadcast withDolbyDigital51(Boolean dolbyDigital51) {
        this.dolbyDigital51 = dolbyDigital51;
        return this;
    }

    public Boolean getDolbySurround() {
        return dolbySurround;
    }

    public void setDolbySurround(Boolean dolbySurround) {
        this.dolbySurround = dolbySurround;
    }

    public ZdfBroadcast withDolbySurround(Boolean dolbySurround) {
        this.dolbySurround = dolbySurround;
        return this;
    }

    public Boolean getDualChannel() {
        return dualChannel;
    }

    public void setDualChannel(Boolean dualChannel) {
        this.dualChannel = dualChannel;
    }

    public ZdfBroadcast withDualChannel(Boolean dualChannel) {
        this.dualChannel = dualChannel;
        return this;
    }

    public Long getDuration() {
        return duration;
    }

    public void setDuration(Long duration) {
        this.duration = duration;
    }

    public ZdfBroadcast withDuration(Long duration) {
        this.duration = duration;
        return this;
    }

    public Long getPartDuration() {
        return partDuration;
    }

    public void setPartDuration(Long partDuration) {
        this.partDuration = partDuration;
    }

    public ZdfBroadcast withPartDuration(Long partDuration) {
        this.partDuration = partDuration;
        return this;
    }

    public Boolean getForeignLangWithCaption() {
        return foreignLangWithCaption;
    }

    public void setForeignLangWithCaption(Boolean foreignLangWithCaption) {
        this.foreignLangWithCaption = foreignLangWithCaption;
    }

    public ZdfBroadcast withForeignLangWithCaption(Boolean foreignLangWithCaption) {
        this.foreignLangWithCaption = foreignLangWithCaption;
        return this;
    }

    public Boolean getHd() {
        return hd;
    }

    public void setHd(Boolean hd) {
        this.hd = hd;
    }

    public ZdfBroadcast withHd(Boolean hd) {
        this.hd = hd;
        return this;
    }

    public String getPosId() {
        return posId;
    }

    public void setPosId(String posId) {
        this.posId = posId;
    }

    public ZdfBroadcast withPosId(String posId) {
        this.posId = posId;
        return this;
    }

    public Long getPartId() {
        return partId;
    }

    public void setPartId(Long partId) {
        this.partId = partId;
    }

    public ZdfBroadcast withPartId(Long partId) {
        this.partId = partId;
        return this;
    }

    public Boolean getLive() {
        return live;
    }

    public void setLive(Boolean live) {
        this.live = live;
    }

    public ZdfBroadcast withLive(Boolean live) {
        this.live = live;
        return this;
    }

    public Boolean getLivestream() {
        return livestream;
    }

    public void setLivestream(Boolean livestream) {
        this.livestream = livestream;
    }

    public ZdfBroadcast withLivestream(Boolean livestream) {
        this.livestream = livestream;
        return this;
    }

    public String getModified() {
        return modified;
    }

    public void setModified(String modified) {
        this.modified = modified;
    }

    public ZdfBroadcast withModified(String modified) {
        this.modified = modified;
        return this;
    }

    public Boolean getMono() {
        return mono;
    }

    public void setMono(Boolean mono) {
        this.mono = mono;
    }

    public ZdfBroadcast withMono(Boolean mono) {
        this.mono = mono;
        return this;
    }

    public Boolean getNewAirtime() {
        return newAirtime;
    }

    public void setNewAirtime(Boolean newAirtime) {
        this.newAirtime = newAirtime;
    }

    public ZdfBroadcast withNewAirtime(Boolean newAirtime) {
        this.newAirtime = newAirtime;
        return this;
    }

    public Boolean getNewProgramData() {
        return newProgramData;
    }

    public void setNewProgramData(Boolean newProgramData) {
        this.newProgramData = newProgramData;
    }

    public ZdfBroadcast withNewProgramData(Boolean newProgramData) {
        this.newProgramData = newProgramData;
        return this;
    }

    public String getPharosId() {
        return pharosId;
    }

    public void setPharosId(String pharosId) {
        this.pharosId = pharosId;
    }

    public ZdfBroadcast withPharosId(String pharosId) {
        this.pharosId = pharosId;
        return this;
    }

    public Boolean getSignLanguage() {
        return signLanguage;
    }

    public void setSignLanguage(Boolean signLanguage) {
        this.signLanguage = signLanguage;
    }

    public ZdfBroadcast withSignLanguage(Boolean signLanguage) {
        this.signLanguage = signLanguage;
        return this;
    }

    public Boolean getStereo() {
        return stereo;
    }

    public void setStereo(Boolean stereo) {
        this.stereo = stereo;
    }

    public ZdfBroadcast withStereo(Boolean stereo) {
        this.stereo = stereo;
        return this;
    }

    public String getVisibleFrom() {
        return visibleFrom;
    }

    public void setVisibleFrom(String visibleFrom) {
        this.visibleFrom = visibleFrom;
    }

    public ZdfBroadcast withVisibleFrom(String visibleFrom) {
        this.visibleFrom = visibleFrom;
        return this;
    }

    public String getVisibleTo() {
        return visibleTo;
    }

    public void setVisibleTo(String visibleTo) {
        this.visibleTo = visibleTo;
    }

    public ZdfBroadcast withVisibleTo(String visibleTo) {
        this.visibleTo = visibleTo;
        return this;
    }

    public String getVpsBegin() {
        return vpsBegin;
    }

    public void setVpsBegin(String vpsBegin) {
        this.vpsBegin = vpsBegin;
    }

    public ZdfBroadcast withVpsBegin(String vpsBegin) {
        this.vpsBegin = vpsBegin;
        return this;
    }

    public Boolean getWidescreen169() {
        return widescreen169;
    }

    public void setWidescreen169(Boolean widescreen169) {
        this.widescreen169 = widescreen169;
    }

    public ZdfBroadcast withWidescreen169(Boolean widescreen169) {
        this.widescreen169 = widescreen169;
        return this;
    }

    public Boolean getWithChat() {
        return withChat;
    }

    public void setWithChat(Boolean withChat) {
        this.withChat = withChat;
    }

    public ZdfBroadcast withWithChat(Boolean withChat) {
        this.withChat = withChat;
        return this;
    }

    public Object getHttpZdfDeRelsCmdmBroadcastsParts() {
        return httpZdfDeRelsCmdmBroadcastsParts;
    }

    public void setHttpZdfDeRelsCmdmBroadcastsParts(Object httpZdfDeRelsCmdmBroadcastsParts) {
        this.httpZdfDeRelsCmdmBroadcastsParts = httpZdfDeRelsCmdmBroadcastsParts;
    }

    public ZdfBroadcast withHttpZdfDeRelsCmdmBroadcastsParts(Object httpZdfDeRelsCmdmBroadcastsParts) {
        this.httpZdfDeRelsCmdmBroadcastsParts = httpZdfDeRelsCmdmBroadcastsParts;
        return this;
    }

    public String getProfile() {
        return profile;
    }

    public void setProfile(String profile) {
        this.profile = profile;
    }

    public ZdfBroadcast withProfile(String profile) {
        this.profile = profile;
        return this;
    }

    public String getSelf() {
        return self;
    }

    public void setSelf(String self) {
        this.self = self;
    }

    public ZdfBroadcast withSelf(String self) {
        this.self = self;
        return this;
    }

    public ZdfTvService getTvService() {
        return zdfTvService;
    }

    public void setTvService(ZdfTvService zdfTvService) {
        this.zdfTvService = zdfTvService;
    }

    public String getProgrammeItem() {
        return httpZdfDeRelsCmdmProgrammeItem;
    }

    public void setHttpZdfDeRelsCmdmProgrammeItem(String httpZdfDeRelsCmdmProgrammeItem) {
        this.httpZdfDeRelsCmdmProgrammeItem = httpZdfDeRelsCmdmProgrammeItem;
    }

    public ZdfBroadcast withHttpZdfDeRelsCmdmProgrammeItem(String httpZdfDeRelsCmdmProgrammeItem) {
        this.httpZdfDeRelsCmdmProgrammeItem = httpZdfDeRelsCmdmProgrammeItem;
        return this;
    }

    public String getSubtitle() {
        return subtitle;
    }

    public void setSubtitle(String subtitle) {
        this.subtitle = subtitle;
    }

    public ZdfBroadcast withSubtitle(String subtitle) {
        this.subtitle = subtitle;
        return this;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    public ZdfBroadcast withText(String text) {
        this.text = text;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public ZdfBroadcast withTitle(String title) {
        this.title = title;
        return this;
    }

    public String getSubheadline() {
        return subheadline;
    }

    public void setSubheadline(String subheadline) {
        this.subheadline = subheadline;
    }

    public ZdfBroadcast withSubheadline(String subheadline) {
        this.subheadline = subheadline;
        return this;
    }

    public String getPrimaryBrandId() {
        return primaryBrandId;
    }

    public void setPrimaryBrandId(String primaryBrandId) {
        this.primaryBrandId = primaryBrandId;
    }

    public ZdfBroadcast withPrimaryBrandId(String primaryBrandId) {
        this.primaryBrandId = primaryBrandId;
        return this;
    }

    public String getPrimaryBrand() {
        return primaryBrand;
    }

    public void setPrimaryBrand(String primaryBrand) {
        this.primaryBrand = primaryBrand;
    }

    public ZdfBroadcast withPrimaryBrand(String primaryBrand) {
        this.primaryBrand = primaryBrand;
        return this;
    }

    public List<String> getBrandIds() {
        return brandIds;
    }

    public void setBrandIds(List<String> brandIds) {
        this.brandIds = brandIds;
    }

    public ZdfBroadcast withBrandIds(List<String> brandIds) {
        this.brandIds = brandIds;
        return this;
    }

    public List<String> getBrandNames() {
        return brandNames;
    }

    public void setBrandNames(List<String> brandNames) {
        this.brandNames = brandNames;
    }

    public ZdfBroadcast withBrandNames(List<String> brandNames) {
        this.brandNames = brandNames;
        return this;
    }

    public ZdfImage getZdfImage() {
        return zdfImage;
    }

    public void setZdfImage(ZdfImage zdfImage) {
        this.zdfImage = zdfImage;
    }

    public ZdfBroadcast withHttpZdfDeRelsImage(ZdfImage zdfImage) {
        this.zdfImage = zdfImage;
        return this;
    }

    public Object getOnlineFrom() {
        return onlineFrom;
    }

    public void setOnlineFrom(Object onlineFrom) {
        this.onlineFrom = onlineFrom;
    }

    public ZdfBroadcast withOnlineFrom(Object onlineFrom) {
        this.onlineFrom = onlineFrom;
        return this;
    }

    public Object getOnlineTo() {
        return onlineTo;
    }

    public void setOnlineTo(Object onlineTo) {
        this.onlineTo = onlineTo;
    }

    public ZdfBroadcast withOnlineTo(Object onlineTo) {
        this.onlineTo = onlineTo;
        return this;
    }

    public Object getGeolocationVOD() {
        return geolocationVOD;
    }

    public void setGeolocationVOD(Object geolocationVOD) {
        this.geolocationVOD = geolocationVOD;
    }

    public ZdfBroadcast withGeolocationVOD(Object geolocationVOD) {
        this.geolocationVOD = geolocationVOD;
        return this;
    }

    public Object getGeolocationLivestream() {
        return geolocationLivestream;
    }

    public void setGeolocationLivestream(Object geolocationLivestream) {
        this.geolocationLivestream = geolocationLivestream;
    }

    public ZdfBroadcast withGeolocationLivestream(Object geolocationLivestream) {
        this.geolocationLivestream = geolocationLivestream;
        return this;
    }

    public Boolean getYoutubeRight() {
        return youtubeRight;
    }

    public void setYoutubeRight(Boolean youtubeRight) {
        this.youtubeRight = youtubeRight;
    }

    public ZdfBroadcast withYoutubeRight(Boolean youtubeRight) {
        this.youtubeRight = youtubeRight;
        return this;
    }

    public Boolean getOnlineFirst() {
        return onlineFirst;
    }

    public void setOnlineFirst(Boolean onlineFirst) {
        this.onlineFirst = onlineFirst;
    }

    public ZdfBroadcast withOnlineFirst(Boolean onlineFirst) {
        this.onlineFirst = onlineFirst;
        return this;
    }

    public Object getHttpZdfDeRelsCmdmSeries() {
        return httpZdfDeRelsCmdmSeries;
    }

    public void setHttpZdfDeRelsCmdmSeries(Object httpZdfDeRelsCmdmSeries) {
        this.httpZdfDeRelsCmdmSeries = httpZdfDeRelsCmdmSeries;
    }

    public ZdfBroadcast withHttpZdfDeRelsCmdmSeries(Object httpZdfDeRelsCmdmSeries) {
        this.httpZdfDeRelsCmdmSeries = httpZdfDeRelsCmdmSeries;
        return this;
    }

    public Object getHttpZdfDeRelsCmdmSeason() {
        return httpZdfDeRelsCmdmSeason;
    }

    public void setHttpZdfDeRelsCmdmSeason(Object httpZdfDeRelsCmdmSeason) {
        this.httpZdfDeRelsCmdmSeason = httpZdfDeRelsCmdmSeason;
    }

    public ZdfBroadcast withHttpZdfDeRelsCmdmSeason(Object httpZdfDeRelsCmdmSeason) {
        this.httpZdfDeRelsCmdmSeason = httpZdfDeRelsCmdmSeason;
        return this;
    }

}
