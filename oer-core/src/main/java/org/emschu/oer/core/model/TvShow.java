package org.emschu.oer.core.model;

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

import javax.persistence.*;
import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Objects;

@SuppressWarnings("squid:S3437")
@Entity
@Table(uniqueConstraints = {
        @UniqueConstraint(columnNames = {"adapterFamily", "technicalId"})
})
public class TvShow implements UpdateControlInterface, Serializable {

    @Id
    @Column(name = "id")
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column
    private LocalDateTime createdAt;

    @Column
    private LocalDateTime updatedAt;

    @Column(length = 1000, nullable = false, name = "title")
    private String title;

    @Column(length = 1500)
    private String url;

    @Column(length = 1500)
    private String homePage;

    @Column(nullable = false, name = "adapterFamily")
    private Channel.AdapterFamily adapterFamily;

    @Column(length = 1500)
    private String imageUrl;

    @Column(length = 32, nullable = false, name = "technicalId")
    private String technicalId;

    @Column(length = 1000, nullable = true, name = "additionalId")
    private String additionalId;

    @ManyToMany(fetch = FetchType.EAGER, targetEntity = ProgramEntry.class)
    @JoinTable(
            joinColumns = @JoinColumn(name = "tv_show_id"),
            inverseJoinColumns = @JoinColumn(name = "related_program_entry_id")
    )
    private List<ProgramEntry> relatedProgramEntries = null;

    @ManyToMany(fetch = FetchType.LAZY, targetEntity = Tag.class)
    @JoinTable(
            joinColumns = @JoinColumn(name = "tv_show_id"),
            inverseJoinColumns = @JoinColumn(name = "tag_id")
    )
    private List<Tag> tags = null;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }

    public LocalDateTime getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(LocalDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public Channel.AdapterFamily getAdapterFamily() {
        return adapterFamily;
    }

    public void setAdapterFamily(Channel.AdapterFamily adapterFamily) {
        this.adapterFamily = adapterFamily;
    }

    public String getImageUrl() {
        return imageUrl;
    }

    public void setImageUrl(String imageUrl) {
        this.imageUrl = imageUrl;
    }

    public String getHomePage() {
        return homePage;
    }

    public void setHomePage(String homePage) {
        this.homePage = homePage;
    }

    public List<ProgramEntry> getRelatedProgramEntries() {
        return relatedProgramEntries;
    }

    public void setRelatedProgramEntries(List<ProgramEntry> relatedProgramEntries) {
        this.relatedProgramEntries = relatedProgramEntries;
    }

    public String getTechnicalId() {
        return technicalId;
    }

    public void setTechnicalId(String technicalId) {
        this.technicalId = technicalId;
    }

    public List<Tag> getTags() {
        return tags;
    }

    public void setTags(List<Tag> tags) {
        this.tags = tags;
    }

    public String getAdditionalId() {
        return additionalId;
    }

    public void setAdditionalId(String additionalId) {
        this.additionalId = additionalId;
    }

    @Override
    public String toString() {
        return "TvShow{" +
                "id=" + id +
                ", createdAt=" + createdAt +
                ", updatedAt=" + updatedAt +
                ", title='" + title + '\'' +
                ", url='" + url + '\'' +
                ", homePage='" + homePage + '\'' +
                ", adapterFamily=" + adapterFamily +
                ", imageUrl='" + imageUrl + '\'' +
                ", technicalId='" + technicalId + '\'' +
                ", additionalId='" + additionalId + '\'' +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TvShow tvShow = (TvShow) o;
        return Objects.equals(id, tvShow.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash("tv-show" + id);
    }
}
