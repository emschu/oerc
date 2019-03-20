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
import java.util.Objects;

@Entity
public class Channel implements Serializable {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private AdapterFamily adapterFamily;
    @Column(unique = true)
    private ChannelKey channelKey;
    @Column(unique = true)
    private String technicalId;
    private String name;
    private String homePage;

    public Channel() {
    }

    /**
     * constructor
     *
     * @param adapterFamily
     * @param channelKey
     * @param technicalId
     * @param name
     * @param homePage
     */
    public Channel(AdapterFamily adapterFamily, ChannelKey channelKey, String technicalId, String name, String homePage) {
        this.adapterFamily = adapterFamily;
        this.channelKey = channelKey;
        this.technicalId = technicalId;
        this.name = name;
        this.homePage = homePage;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public AdapterFamily getAdapterFamily() {
        return adapterFamily;
    }

    public void setAdapterFamily(AdapterFamily adapterFamily) {
        this.adapterFamily = adapterFamily;
    }

    public String getTechnicalId() {
        return technicalId;
    }

    public void setTechnicalId(String technicalId) {
        this.technicalId = technicalId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getHomePage() {
        return homePage;
    }

    public void setHomePage(String homePage) {
        this.homePage = homePage;
    }

    public ChannelKey getChannelKey() {
        return channelKey;
    }

    public void setChannelKey(ChannelKey channelKey) {
        this.channelKey = channelKey;
    }

    /* --- enums --- */
    public enum AdapterFamily {
        ARD("ARD"), ZDF("ZDF");

        private String familyKey;
        AdapterFamily(String familyKey) {
            this.familyKey = familyKey;
        }

        @Override
        public String toString() {
            return familyKey;
        }
    }

    public enum ChannelKey {
        // ARD tv sender family:
        ARD, ZDF, ZDF_INFO, ZDF_NEO, DREISAT, ARTE, BR, HR, MDR, NDR, RBB, RADIO_BREMEN_TV, SR, SWR_BW, SWR_RP, WDR, ALPHA, TAGESSCHAU_24, ONE, KIKA, PHOENIX
    }

    @Override
    public String toString() {
        return "Channel{" +
                "id=" + id +
                ", adapterFamily=" + adapterFamily +
                ", channelKey=" + channelKey +
                ", technicalId='" + technicalId + '\'' +
                ", name='" + name + '\'' +
                ", homePage='" + homePage + '\'' +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Channel channel = (Channel) o;
        return Objects.equals(id, channel.id) &&
                adapterFamily == channel.adapterFamily &&
                channelKey == channel.channelKey &&
                Objects.equals(technicalId, channel.technicalId) &&
                Objects.equals(name, channel.name);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id, adapterFamily, channelKey, technicalId, name);
    }
}
