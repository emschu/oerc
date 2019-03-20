
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

import java.util.Iterator;
import java.util.List;
import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

public class ProgramItemModel implements Iterable <ZdfBroadcast> {

    @SerializedName("profile")
    @Expose
    private String profile;
    @SerializedName("http://zdf.de/rels/cmdm/broadcasts")
    @Expose
    private List<ZdfBroadcast> zdfBroadcasts = null;
    @SerializedName("self")
    @Expose
    private String self;

    public String getProfile() {
        return profile;
    }

    public void setProfile(String profile) {
        this.profile = profile;
    }

    public List<ZdfBroadcast> getZdfBroadcasts() {
        return zdfBroadcasts;
    }

    public void setZdfBroadcasts(List<ZdfBroadcast> zdfBroadcasts) {
        this.zdfBroadcasts = zdfBroadcasts;
    }

    public String getSelf() {
        return self;
    }

    public void setSelf(String self) {
        this.self = self;
    }

    @Override
    public Iterator<ZdfBroadcast> iterator() {
        return zdfBroadcasts.iterator();
    }
}
