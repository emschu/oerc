
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

public class ZdfTvShowResponseModel extends AbstractZdfApiElement {

    @SerializedName("id")
    @Expose
    private String id;
    @SerializedName("externalId")
    @Expose
    private String externalId;
    @SerializedName("contentType")
    @Expose
    private String contentType;
    @SerializedName("publicationDate")
    @Expose
    private String publicationDate;
    @SerializedName("modificationDate")
    @Expose
    private String modificationDate;
    @SerializedName("title")
    @Expose
    private String title;
    @SerializedName("editorialDate")
    @Expose
    private String editorialDate;
    @SerializedName("teaserImageRef")
    @Expose
    private TeaserImageRef teaserImageRef;

    @SerializedName("hasVideo")
    @Expose
    private Boolean hasVideo;
    @SerializedName("hideFacebookShareButton")
    @Expose
    private Boolean hideFacebookShareButton;
    @SerializedName("hideGoogleShareButton")
    @Expose
    private Boolean hideGoogleShareButton;
    @SerializedName("hideTwitterShareButton")
    @Expose
    private Boolean hideTwitterShareButton;
    @SerializedName("structureNodePath")
    @Expose
    private String structureNodePath;
    @SerializedName("http://zdf.de/rels/category")
    @Expose
    private HttpZdfDeRelsCategory httpZdfDeRelsCategory;
    @SerializedName("http://zdf.de/rels/uri")
    @Expose
    private String httpZdfDeRelsUri;
    @SerializedName("http://zdf.de/rels/sharing-url")
    @Expose
    private String httpZdfDeRelsSharingUrl;
    @SerializedName("http://zdf.de/rels/brand")
    @Expose
    private HttpZdfDeRelsBrand httpZdfDeRelsBrand;
    @SerializedName("http://zdf.de/rels/docstats")
    @Expose
    private HttpZdfDeRelsDocstats httpZdfDeRelsDocstats;
    @SerializedName("http://zdf.de/rels/rss/feed")
    @Expose
    private String httpZdfDeRelsRssFeed;
    @SerializedName("stage")
    @Expose
    private List<Stage> stage = new ArrayList<Stage>();
    @SerializedName("editorial-tag")
    @Expose
    private List<EditorialTag> editorialTag = new ArrayList<EditorialTag>();
    @SerializedName("module")
    @Expose
    private List<Module> module = new ArrayList<Module>();
    @SerializedName("http://zdf.de/rels/content/conf-section")
    @Expose
    private HttpZdfDeRelsContentConfSection httpZdfDeRelsContentConfSection;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getExternalId() {
        return externalId;
    }

    public void setExternalId(String externalId) {
        this.externalId = externalId;
    }

    public String getContentType() {
        return contentType;
    }

    public void setContentType(String contentType) {
        this.contentType = contentType;
    }

    public String getPublicationDate() {
        return publicationDate;
    }

    public void setPublicationDate(String publicationDate) {
        this.publicationDate = publicationDate;
    }

    public String getModificationDate() {
        return modificationDate;
    }

    public void setModificationDate(String modificationDate) {
        this.modificationDate = modificationDate;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getEditorialDate() {
        return editorialDate;
    }

    public void setEditorialDate(String editorialDate) {
        this.editorialDate = editorialDate;
    }

    public TeaserImageRef getTeaserImageRef() {
        return teaserImageRef;
    }

    public void setTeaserImageRef(TeaserImageRef teaserImageRef) {
        this.teaserImageRef = teaserImageRef;
    }

    public Boolean getHasVideo() {
        return hasVideo;
    }

    public void setHasVideo(Boolean hasVideo) {
        this.hasVideo = hasVideo;
    }

    public Boolean getHideFacebookShareButton() {
        return hideFacebookShareButton;
    }

    public void setHideFacebookShareButton(Boolean hideFacebookShareButton) {
        this.hideFacebookShareButton = hideFacebookShareButton;
    }

    public Boolean getHideGoogleShareButton() {
        return hideGoogleShareButton;
    }

    public void setHideGoogleShareButton(Boolean hideGoogleShareButton) {
        this.hideGoogleShareButton = hideGoogleShareButton;
    }

    public Boolean getHideTwitterShareButton() {
        return hideTwitterShareButton;
    }

    public void setHideTwitterShareButton(Boolean hideTwitterShareButton) {
        this.hideTwitterShareButton = hideTwitterShareButton;
    }

    public String getStructureNodePath() {
        return structureNodePath;
    }

    public void setStructureNodePath(String structureNodePath) {
        this.structureNodePath = structureNodePath;
    }

    public HttpZdfDeRelsCategory getHttpZdfDeRelsCategory() {
        return httpZdfDeRelsCategory;
    }

    public void setHttpZdfDeRelsCategory(HttpZdfDeRelsCategory httpZdfDeRelsCategory) {
        this.httpZdfDeRelsCategory = httpZdfDeRelsCategory;
    }

    public String getHttpZdfDeRelsUri() {
        return httpZdfDeRelsUri;
    }

    public void setHttpZdfDeRelsUri(String httpZdfDeRelsUri) {
        this.httpZdfDeRelsUri = httpZdfDeRelsUri;
    }

    public String getHttpZdfDeRelsSharingUrl() {
        return httpZdfDeRelsSharingUrl;
    }

    public void setHttpZdfDeRelsSharingUrl(String httpZdfDeRelsSharingUrl) {
        this.httpZdfDeRelsSharingUrl = httpZdfDeRelsSharingUrl;
    }

    public HttpZdfDeRelsBrand getHttpZdfDeRelsBrand() {
        return httpZdfDeRelsBrand;
    }

    public void setHttpZdfDeRelsBrand(HttpZdfDeRelsBrand httpZdfDeRelsBrand) {
        this.httpZdfDeRelsBrand = httpZdfDeRelsBrand;
    }

    public HttpZdfDeRelsDocstats getHttpZdfDeRelsDocstats() {
        return httpZdfDeRelsDocstats;
    }

    public void setHttpZdfDeRelsDocstats(HttpZdfDeRelsDocstats httpZdfDeRelsDocstats) {
        this.httpZdfDeRelsDocstats = httpZdfDeRelsDocstats;
    }

    public String getHttpZdfDeRelsRssFeed() {
        return httpZdfDeRelsRssFeed;
    }

    public void setHttpZdfDeRelsRssFeed(String httpZdfDeRelsRssFeed) {
        this.httpZdfDeRelsRssFeed = httpZdfDeRelsRssFeed;
    }

    public List<Stage> getStage() {
        return stage;
    }

    public void setStage(List<Stage> stage) {
        this.stage = stage;
    }

    public List<EditorialTag> getEditorialTag() {
        return editorialTag;
    }

    public void setEditorialTag(List<EditorialTag> editorialTag) {
        this.editorialTag = editorialTag;
    }

    public List<Module> getModule() {
        return module;
    }

    public void setModule(List<Module> module) {
        this.module = module;
    }

    public HttpZdfDeRelsContentConfSection getHttpZdfDeRelsContentConfSection() {
        return httpZdfDeRelsContentConfSection;
    }

    public void setHttpZdfDeRelsContentConfSection(HttpZdfDeRelsContentConfSection httpZdfDeRelsContentConfSection) {
        this.httpZdfDeRelsContentConfSection = httpZdfDeRelsContentConfSection;
    }
}
