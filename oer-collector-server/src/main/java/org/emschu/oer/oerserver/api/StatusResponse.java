package org.emschu.oer.oerserver.api;

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

import org.emschu.oer.core.model.Channel;

import java.time.LocalDateTime;
import java.util.List;

public class StatusResponse {
    // counter
    private long artistCount;
    private long channelCount;
    private long imageLinksCount;
    private long programEntryCount;
    private long tagCount;
    private long tvShowCount;

    // api information
    private boolean isCurrentlyUpdating;
    private String version;
    private LocalDateTime serverDateTime;

    // tv channels
    private List<Channel> tvChannels;

    public long getArtistCount() {
        return artistCount;
    }

    public void setArtistCount(long artistCount) {
        this.artistCount = artistCount;
    }

    public long getChannelCount() {
        return channelCount;
    }

    public void setChannelCount(long channelCount) {
        this.channelCount = channelCount;
    }

    public long getImageLinksCount() {
        return imageLinksCount;
    }

    public void setImageLinksCount(long imageLinksCount) {
        this.imageLinksCount = imageLinksCount;
    }

    public long getProgramEntryCount() {
        return programEntryCount;
    }

    public void setProgramEntryCount(long programEntryCount) {
        this.programEntryCount = programEntryCount;
    }

    public long getTagCount() {
        return tagCount;
    }

    public void setTagCount(long tagCount) {
        this.tagCount = tagCount;
    }

    public long getTvShowCount() {
        return tvShowCount;
    }

    public void setTvShowCount(long tvShowCount) {
        this.tvShowCount = tvShowCount;
    }

    public List<Channel> getTvChannels() {
        return tvChannels;
    }

    public void setTvChannels(List<Channel> tvChannels) {
        this.tvChannels = tvChannels;
    }

    public boolean isCurrentlyUpdating() {
        return isCurrentlyUpdating;
    }

    public void setCurrentlyUpdating(boolean currentlyUpdating) {
        isCurrentlyUpdating = currentlyUpdating;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public LocalDateTime getServerDateTime() {
        return serverDateTime;
    }

    public void setServerDateTime(LocalDateTime serverDateTime) {
        this.serverDateTime = serverDateTime;
    }

}
