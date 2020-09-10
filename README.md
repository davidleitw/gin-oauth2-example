# gin-oauth2-example

## a simple example use gin framework to login with oauth2. 

**此專案基於 golang.org/x/oauth2 來實作一個簡單的登入系統，串接google, facebook, 以及github。** 

**[Demo](https://ginoauth-example.herokuapp.com/login) 網址，僅要求最低權限，提供email以及username。**

以下會依照個人理解以及網路上心得文章的整理， 做一些簡單的觀念介紹。

- # Oauth

  - ## 使用場景，為什麼我們要使用Oauth2?
    假設今天我們要開發一個應用程式，功能是可以將使用者存在google帳號底下的照片都匯入應用程式並且做一些處理，這時候就要想說，我們的應用程式要怎麼存取到使用者帳號底下的資源呢？  
    <br>
    
    在以前我們可能會想說，使用者自行輸入他的google帳號密碼，我們開發的應用程式再藉由那組帳號密碼獲得存取照片的權限，藉此來對照片進行處理。 但是這樣會衍生出很多問題，主要就是如果該應用程式是藉由我們使用者的google帳號密碼來獲得權限的話，應用程式本身獲得的權限太多了，應用程式本來只需要照片的存取權而已，如果是拿google帳密的話則是拿到了整個帳號的使用權限。  
    <br>
    
    如果現在有個不懷好意的應用程式獲得了你整個google帳號密碼的使用權限， 便可以透過它去登入很多不同的網站，又或者是拿去做非法的交易等等...  所以需要提出一個解決方法，一來是可以讓應用程式獲得他想要獲得的資源，二來是也不會開放太多的權限給應用程 式，只授權部份的存取功能。   
    <br>
    
    從另外一個角度來討論， 如果今天我很信任這個應用程式，把整個帳號的權限都交給了這個應
    用程式， 某一天決定之後可能不再使用這個應用程式了。 這時候我們如果想要把權限收回來怎麼做呢?  我們需要把帳號密碼重新設定一輪，才可以確保應用程式本身不再有獲得我帳號權限的能力。  


    ### 總結一下傳統方法的缺點:
    (1) 應用程式為了後續的操作可能會保留使用者的帳號密碼， 我們無法確定是否安全。
    (2) 使用者無法限制應用程式能掌握到的權限，以及權限有效的時間。
    (3) 修改密碼之後可以收回權限，不過也同時讓所有獲得用戶授權的第三方應用全部失效。
    (4) 只要有一個第三方的應用程式出現了資料外洩問題，會導致其他第三方應用也有被盜用的危險。

    </br>
        
    > **簡單的說Oauth就是一種授權機制，以google為存放data的地方為例， 使用者告知google，授權第三方的應用程式可以獲得部份的資源， google的系統則產生一個token， 第三方的應用可以透過token來在實現內進行授權權限內的資料存取。**


    > **Token vs PassWord**
    > **(1) token是有時效性，時間到了會自動撤銷，不會讓第三方應用程式一直持有權限。**
    > **(2) 如果今天不想要繼續授權給應用程式了，資源擁有者可以隨時撤銷token的有效性。**
    > **(3) token最大的好處是有權限的管理，可以只開放部份的資源存取權給第三方的應用程式。**

  - ## Oauth2.0標準文件 [RFC6749](https://tools.ietf.org/html/rfc6749)
    這份規格書內清楚的寫了Oauth2.0的設計準則， 如果想要清楚的知道oauth2.0底層是如何運作的可以參考看看。
    值得注意的是在RFC 6749的文件內清楚的說明了Oauth2.0的角色
    > **Oauth在傳統的架構上引入了一層授權層，用以分隔客戶端以及資源擁有者，當資源的擁有者授權客戶端(第三方應用程式)應用可以存取資源的時候，用來存放資源的伺服器會頒發一個Access Token給客戶端(第三方應用程式)， 客戶端拿到Token之後就可以藉由此Token對資源做有限度的存取(並沒有拿到全部的權限)。**

    > 原文: 
    > The OAuth 2.0 authorization framework enables a third-party   application to obtain limited access to an HTTP service, either on   behalf of a resource owner by orchestrating an approval interaction   between the resource owner and the HTTP service, or by allowing the   third-party application to obtain access on its own behalf.  This   specification replaces and obsoletes the OAuth 1.0 protocol described   in RFC 5849.
    
    </br>

    由上面那段重點可以得知Oauth扮演的核心角色就是向第三方應用程式頒發Token，作為第三方應用跟資源網站的橋樑。

  - ## Authorization Grant(四種獲得Token的方式)
    > #### 1. **authorization-code (授權碼)**
    > #### 2. **implicit (隱藏式)**
    > #### 3. **password (密碼式)**
    > #### 4. **client credentials (客戶端憑證)**
    
    接著就依序介紹一下這四種獲得Token的方式個別用在哪種場景，以及需要注意的地方。
    </br>
    
    - ### Authorization-Code(授權碼):
        此方法是目前最為常見的一種手法，在主流的前後端分離架構上，通常採取這種方法，使得可以在**後端獲得Token**，所有有關於資源存取的運算都放在後端， 可以減少token暴露的機會。




