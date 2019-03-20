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

public class TvShowModel {

    @SerializedName("actorDetails")
    @Expose
    private Object actorDetails;
    @SerializedName("category")
    @Expose
    private String category;
    @SerializedName("country")
    @Expose
    private Object country;
    @SerializedName("crewDetails")
    @Expose
    private CrewDetails crewDetails;
    @SerializedName("genre")
    @Expose
    private String genre;
    @SerializedName("guest")
    @Expose
    private List<Object> guest = null;
    @SerializedName("language")
    @Expose
    private String language;
    @SerializedName("originalTitle")
    @Expose
    private Object originalTitle;
    @SerializedName("subtitle")
    @Expose
    private Object subtitle;
    @SerializedName("text")
    @Expose
    private Object text;
    @SerializedName("title")
    @Expose
    private String title;
    @SerializedName("year")
    @Expose
    private Object year;
    @SerializedName("updateTimestamp")
    @Expose
    private String updateTimestamp;
    @SerializedName("contentId")
    @Expose
    private String contentId;
    @SerializedName("profile")
    @Expose
    private String profile;
    @SerializedName("self")
    @Expose
    private String self;
    @SerializedName("textVariant")
    @Expose
    private Long textVariant;
    @SerializedName("textShort")
    @Expose
    private Object textShort;
    @SerializedName("textShortCont")
    @Expose
    private Object textShortCont;
    @SerializedName("textCont")
    @Expose
    private Object textCont;
    @SerializedName("textAccomp")
    @Expose
    private Object textAccomp;
    @SerializedName("textLocation")
    @Expose
    private Object textLocation;
    @SerializedName("textPresentation")
    @Expose
    private Object textPresentation;
    @SerializedName("textReporter")
    @Expose
    private Object textReporter;
    @SerializedName("textExpert")
    @Expose
    private Object textExpert;
    @SerializedName("textFree")
    @Expose
    private Object textFree;
    @SerializedName("textOrder")
    @Expose
    private Object textOrder;
    @SerializedName("fsk")
    @Expose
    private Object fsk;
    @SerializedName("brandIds")
    @Expose
    private List<String> brandIds = null;
    @SerializedName("brandNames")
    @Expose
    private List<String> brandNames = null;
    @SerializedName("hinttext")
    @Expose
    private Object hinttext;
    @SerializedName("interactiv")
    @Expose
    private Boolean interactiv;
    @SerializedName("subheadline")
    @Expose
    private Object subheadline;
    @SerializedName("primaryBrandId")
    @Expose
    private String primaryBrandId;
    @SerializedName("primaryBrand")
    @Expose
    private String primaryBrand;
    @SerializedName("searchable")
    @Expose
    private Boolean searchable;
    @SerializedName("tvTipp")
    @Expose
    private Boolean tvTipp;
    @SerializedName("tvTippDetail")
    @Expose
    private Object tvTippDetail;
    @SerializedName("videoOnDemand")
    @Expose
    private Boolean videoOnDemand;
    @SerializedName("http://zdf.de/rels/image")
    @Expose
    private ZdfImage zdfImage;
    @SerializedName("http://zdf.de/rels/content/video-page")
    @Expose
    private String httpZdfDeRelsContentVideoPage;
    @SerializedName("http://zdf.de/rels/cmdm/broadcasts")
    @Expose
    private List<ZdfBroadcast> broadCasts = null;
    @SerializedName("episodeNumber")
    @Expose
    private Object episodeNumber;
    @SerializedName("episodeImdbId")
    @Expose
    private Object episodeImdbId;
    @SerializedName("originalAirDate")
    @Expose
    private Object originalAirDate;
    @SerializedName("http://zdf.de/rels/cmdm/brand")
    @Expose
    private Object httpZdfDeRelsCmdmBrand;
    @SerializedName("http://zdf.de/rels/cmdm/series")
    @Expose
    private Object httpZdfDeRelsCmdmSeries;
    @SerializedName("http://zdf.de/rels/cmdm/season")
    @Expose
    private Object httpZdfDeRelsCmdmSeason;

    public Object getActorDetails() {
        return actorDetails;
    }

    public void setActorDetails(Object actorDetails) {
        this.actorDetails = actorDetails;
    }

    public String getCategory() {
        return category;
    }

    public void setCategory(String category) {
        this.category = category;
    }

    public Object getCountry() {
        return country;
    }

    public void setCountry(Object country) {
        this.country = country;
    }

    public CrewDetails getCrewDetails() {
        return crewDetails;
    }

    public void setCrewDetails(CrewDetails crewDetails) {
        this.crewDetails = crewDetails;
    }

    public String getGenre() {
        return genre;
    }

    public void setGenre(String genre) {
        this.genre = genre;
    }

    public List<Object> getGuest() {
        return guest;
    }

    public void setGuest(List<Object> guest) {
        this.guest = guest;
    }

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }

    public Object getOriginalTitle() {
        return originalTitle;
    }

    public void setOriginalTitle(Object originalTitle) {
        this.originalTitle = originalTitle;
    }

    public Object getSubtitle() {
        return subtitle;
    }

    public void setSubtitle(Object subtitle) {
        this.subtitle = subtitle;
    }

    public Object getText() {
        return text;
    }

    public void setText(Object text) {
        this.text = text;
    }

    public TvShowModel withText(Object text) {
        this.text = text;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public Object getYear() {
        return year;
    }

    public void setYear(Object year) {
        this.year = year;
    }

    public String getUpdateTimestamp() {
        return updateTimestamp;
    }

    public void setUpdateTimestamp(String updateTimestamp) {
        this.updateTimestamp = updateTimestamp;
    }

    public String getContentId() {
        return contentId;
    }

    public void setContentId(String contentId) {
        this.contentId = contentId;
    }

    public String getProfile() {
        return profile;
    }

    public void setProfile(String profile) {
        this.profile = profile;
    }

    public String getSelf() {
        return self;
    }

    public void setSelf(String self) {
        this.self = self;
    }

    public Long getTextVariant() {
        return textVariant;
    }

    public void setTextVariant(Long textVariant) {
        this.textVariant = textVariant;
    }

    public Object getTextShort() {
        return textShort;
    }

    public void setTextShort(Object textShort) {
        this.textShort = textShort;
    }

    public Object getTextShortCont() {
        return textShortCont;
    }

    public void setTextShortCont(Object textShortCont) {
        this.textShortCont = textShortCont;
    }

    public Object getTextCont() {
        return textCont;
    }

    public void setTextCont(Object textCont) {
        this.textCont = textCont;
    }

    public Object getTextAccomp() {
        return textAccomp;
    }

    public void setTextAccomp(Object textAccomp) {
        this.textAccomp = textAccomp;
    }

    public Object getTextLocation() {
        return textLocation;
    }

    public void setTextLocation(Object textLocation) {
        this.textLocation = textLocation;
    }

    public Object getTextPresentation() {
        return textPresentation;
    }

    public void setTextPresentation(Object textPresentation) {
        this.textPresentation = textPresentation;
    }

    public Object getTextReporter() {
        return textReporter;
    }

    public void setTextReporter(Object textReporter) {
        this.textReporter = textReporter;
    }

    public Object getTextExpert() {
        return textExpert;
    }

    public void setTextExpert(Object textExpert) {
        this.textExpert = textExpert;
    }

    public Object getTextFree() {
        return textFree;
    }

    public void setTextFree(Object textFree) {
        this.textFree = textFree;
    }

    public Object getTextOrder() {
        return textOrder;
    }

    public void setTextOrder(Object textOrder) {
        this.textOrder = textOrder;
    }

    public Object getFsk() {
        return fsk;
    }

    public void setFsk(Object fsk) {
        this.fsk = fsk;
    }

    public List<String> getBrandIds() {
        return brandIds;
    }

    public void setBrandIds(List<String> brandIds) {
        this.brandIds = brandIds;
    }

    public List<String> getBrandNames() {
        return brandNames;
    }

    public void setBrandNames(List<String> brandNames) {
        this.brandNames = brandNames;
    }

    public Object getHinttext() {
        return hinttext;
    }

    public void setHinttext(Object hinttext) {
        this.hinttext = hinttext;
    }

    public Boolean getInteractiv() {
        return interactiv;
    }

    public void setInteractiv(Boolean interactiv) {
        this.interactiv = interactiv;
    }

    public Object getSubheadline() {
        return subheadline;
    }

    public void setSubheadline(Object subheadline) {
        this.subheadline = subheadline;
    }

    public String getPrimaryBrandId() {
        return primaryBrandId;
    }

    public void setPrimaryBrandId(String primaryBrandId) {
        this.primaryBrandId = primaryBrandId;
    }

    public String getPrimaryBrand() {
        return primaryBrand;
    }

    public void setPrimaryBrand(String primaryBrand) {
        this.primaryBrand = primaryBrand;
    }

    public Boolean getSearchable() {
        return searchable;
    }

    public void setSearchable(Boolean searchable) {
        this.searchable = searchable;
    }

    public Boolean getTvTipp() {
        return tvTipp;
    }

    public void setTvTipp(Boolean tvTipp) {
        this.tvTipp = tvTipp;
    }

    public Object getTvTippDetail() {
        return tvTippDetail;
    }

    public void setTvTippDetail(Object tvTippDetail) {
        this.tvTippDetail = tvTippDetail;
    }

    public Boolean getVideoOnDemand() {
        return videoOnDemand;
    }

    public void setVideoOnDemand(Boolean videoOnDemand) {
        this.videoOnDemand = videoOnDemand;
    }

    public String getHttpZdfDeRelsContentVideoPage() {
        return httpZdfDeRelsContentVideoPage;
    }

    public void setHttpZdfDeRelsContentVideoPage(String httpZdfDeRelsContentVideoPage) {
        this.httpZdfDeRelsContentVideoPage = httpZdfDeRelsContentVideoPage;
    }

    public Object getEpisodeNumber() {
        return episodeNumber;
    }

    public void setEpisodeNumber(Object episodeNumber) {
        this.episodeNumber = episodeNumber;
    }

    public Object getEpisodeImdbId() {
        return episodeImdbId;
    }

    public void setEpisodeImdbId(Object episodeImdbId) {
        this.episodeImdbId = episodeImdbId;
    }

    public Object getOriginalAirDate() {
        return originalAirDate;
    }

    public void setOriginalAirDate(Object originalAirDate) {
        this.originalAirDate = originalAirDate;
    }

    public Object getHttpZdfDeRelsCmdmBrand() {
        return httpZdfDeRelsCmdmBrand;
    }

    public void setHttpZdfDeRelsCmdmBrand(Object httpZdfDeRelsCmdmBrand) {
        this.httpZdfDeRelsCmdmBrand = httpZdfDeRelsCmdmBrand;
    }

    public Object getHttpZdfDeRelsCmdmSeries() {
        return httpZdfDeRelsCmdmSeries;
    }

    public void setHttpZdfDeRelsCmdmSeries(Object httpZdfDeRelsCmdmSeries) {
        this.httpZdfDeRelsCmdmSeries = httpZdfDeRelsCmdmSeries;
    }

    public Object getHttpZdfDeRelsCmdmSeason() {
        return httpZdfDeRelsCmdmSeason;
    }

    public void setHttpZdfDeRelsCmdmSeason(Object httpZdfDeRelsCmdmSeason) {
        this.httpZdfDeRelsCmdmSeason = httpZdfDeRelsCmdmSeason;
    }
}
