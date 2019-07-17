package org.emschu.oer.collector.reader;

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

import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;
import org.emschu.oer.zdf_api.model.ProgramItemModel;
import org.emschu.oer.zdf_api.model.TvShowModel;
import org.emschu.oer.zdf_api.model.ZdfTvShowResponseModel;
import org.springframework.http.HttpMethod;
import org.springframework.http.client.ClientHttpRequest;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.http.client.SimpleClientHttpRequestFactory;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.*;
import java.util.HashMap;
import java.util.Map;

/**
 * this class wraps all zdf api calls and provides methods to retrieve responses
 */
public class ZdfApiFetcher extends Fetcher {
    private static String CURRENT_ZDF_API_KEY = null;
    private static Map<String, String> headerMap = new HashMap<>();
    private static final SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();

    static {
        headerMap.put("Host", "api.zdf.de");
        headerMap.put("Accept", "application/vnd.de.zdf.v1.0+json");
        headerMap.put("Origin", "https://www.zdf.de");
        headerMap.put("Api-Auth", "");
    }

    public static ProgramItemModel getProgram(String url) {
        return (ProgramItemModel) getDocument(url, ProgramItemModel.class);
    }

    public static TvShowModel getTvShow(String url) {
        return (TvShowModel) getDocument(url, TvShowModel.class);
    }

    public static ZdfTvShowResponseModel getSingleTvShow(String url) {
        return (ZdfTvShowResponseModel) getDocument(url, ZdfTvShowResponseModel.class);
    }

    /**
     * raw url
     *
     * @param url
     * @param outputClass
     * @return has the type defined of outputClass
     */
    @SuppressWarnings("unchecked")
    private static Object getDocument(String url, Class outputClass) {
        if (getCurrentZdfApiKey() == null || url == null || url.isEmpty() || outputClass == null) {
            throw new IllegalStateException("no zdf api key set. can not process api call.");
        }

        if (getProxyHost() != null && getProxyPort() != 0) {
            Proxy proxy = new Proxy(Proxy.Type.HTTP, new InetSocketAddress(getProxyHost(), getProxyPort()));
            requestFactory.setProxy(proxy);
        }

        BufferedReader vr = null;
        try {
            URL newUrl = new URL(url);
            ClientHttpRequest request = requestFactory.createRequest(newUrl.toURI(), HttpMethod.GET);
            request.getHeaders().clear();
            headerMap.forEach((key, value) -> {
                if (key.equals("Api-Auth")) {
                    // api key value is not available during static setup of this class
                    request.getHeaders().add(key, "Bearer " + getCurrentZdfApiKey());
                } else {
                    request.getHeaders().add(key, value);
                }
            });
            ClientHttpResponse response = request.execute();

            if (!response.getStatusCode().is2xxSuccessful()) {
                LOG.warning("Api return status code: " + response.getStatusCode().toString());
                LOG.warning("Request failed for url: " + url);
            }

            increaseCounter();
            vr = new BufferedReader(new InputStreamReader(response.getBody()));
            return new Gson().fromJson(vr, outputClass);
        } catch (IOException | URISyntaxException | JsonSyntaxException e) {
            LOG.warning(String.format("problem fetching '%s'. Throwing exception '%s' with cause '%s'", url, e.getMessage(), e.getCause()));
            LOG.throwing(ZdfApiFetcher.class.getName(), "getProgram", e);
        } finally {
            if (vr != null) {
                try {
                    vr.close();
                } catch (IOException e) {
                    LOG.throwing(ZdfApiFetcher.class.getName(), "getProgram", e);
                    LOG.warning("io exception during close of buffered reader");
                }
            }
        }
        return null;
    }

    public static String getCurrentZdfApiKey() {
        return CURRENT_ZDF_API_KEY;
    }

    public static void setCurrentZdfApiKey(String currentZdfApiKey) {
        CURRENT_ZDF_API_KEY = currentZdfApiKey;
    }
}
