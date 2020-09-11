# gin-oauth2-example

## a simple example use gin framework to login with oauth2. 

**此專案基於 golang.org/x/oauth2 來實作一個簡單的第三方登入系統，串接google, facebook, 以及github。** 

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
    
    從另外一個角度來討論， 如果今天我很信任這個應用程式，把整個帳號的權限都交給了這個應用程式， 某一天決定之後可能不再使用這個應用程式了。 這時候我們如果想要把權限收回來怎麼做呢?  我們需要把帳號密碼重新設定一輪，才可以確保應用程式本身不再有獲得我帳號權限的能力。  


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

  - ## Oauth2.0規格書 [RFC6749](https://tools.ietf.org/html/rfc6749)
    這份規格書內清楚的寫了Oauth2.0的設計準則， 如果想要清楚的知道oauth2.0底層是如何運作的可以參考看看。
    值得注意的是在RFC 6749的文件內清楚的說明了Oauth2.0的角色
    > **Oauth在傳統的架構上引入了一層授權層，用以分隔客戶端以及資源擁有者，當資源的擁有者授權客戶端(第三方應用程式)應用可以存取資源的時候，用來存放資源的伺服器會頒發一個Access Token給客戶端(第三方應用程式)， 客戶端拿到Token之後就可以藉由此Token對資源做有限度的存取(並沒有拿到全部的權限)。**

    > 原文段落:  
    > 
    > The OAuth 2.0 authorization framework enables a third-party   application to obtain limited access to an HTTP service, either on   behalf of a resource owner by orchestrating an approval interaction   between the resource owner and the HTTP service, or by allowing the   third-party application to obtain access on its own behalf.  This   specification replaces and obsoletes the OAuth 1.0 protocol described   in RFC 5849.
    
    </br>

    由上面那段重點可以得知Oauth扮演的核心角色就是向第三方應用程式頒發Token，作為第三方應用跟資源網站的橋樑。
    <br>

    >  **補充: [RFC6750](https://tools.ietf.org/html/rfc6750)  該文件敘述一些有關於Token實作上的細節**

  - ## Authorization Grant(四種獲得Token的方式)
    > #### 1. **authorization-code (授權碼)**
    > #### 2. **implicit (隱藏式)**
    > #### 3. **password (密碼式)**
    > #### 4. **client credentials (客戶端憑證)**
    
    接著就依序介紹一下這四種獲得Token的方式個別用在哪種場景，以及需要注意的地方。
    </br>
    
    - ### Authorization-Code(授權碼):
        **此範例就是使用Authorization Code的方法去實行的。**
        <br>

        此方法是目前最為常見的一種手法，在主流的前後端分離架構上，通常採取這種方法，使得可以在**後端獲得Token**，所有有關於資源存取的運算都放在後端， 可以減少token暴露的機會。
        <br>

        通常如果是使用這種驗證方法的話，需要先去跟資源擁有者申請**ClientID**跟**ClientSecret**。  
        藉由ClientID，資源擁有者(ex: google)才知道今天是哪個第三方應用在請求資源。
        <br>
        
        大致流程如下圖所示 [來源](https://itnext.io/an-oauth-2-0-introduction-for-beginners-6e386b19f7a9)
        ![](https://miro.medium.com/max/3553/1*anmFPvD_EVMiZOo-W76qyA.png)
        <br>

        以下解說皆以此專案Demo的網址作為範例， 可以與[Demo](https://ginoauth-example.herokuapp.com/login)一起服用。  
        一開始我們點進去的網址是**https://ginoauth-example.herokuapp.com/login**

        今天用戶想要以google帳號登入我們的第三方應用(A)， 所以點了google登入的按鍵。

        按下去之後會跳轉到以下網址(範例，僅列出常見的參數)
        > https://accounts.google.com/o/oauth2/auth/identifier?
        > client_id=xxx&
        > response_type=code&
        > redirect_url=https://ginoauth-example.herokuapp.com/callback/google&
        > scope=https://www.googleapis.com/auth/userinfo.profile&
        > state=xxxx 
        
        參數說明
        > - cliend_id: 填入第三方程式申請oauth服務時獲得的ID，主要目的是讓google知道是誰在申請。
        > - response_type: 參數代表目前申請token採用的是哪一種方式， 這邊填入"code"，代表要申請授權碼。
        > - redirect_url: 參數代表如果google接受請求之後會跳轉的界面。
        > - scope: 參數是代表這次第三方應用申請oauth之後供存取的權限(授權範圍)。
        > - state: 防止CSRF攻擊

        跳轉後由user輸入帳號密碼，告知google授權資源存取權給A。(授權的範圍就是scope中所指定的資源)  
        google確認user同意之後會將網頁跳轉到上方redirect_url參數中所指定的url。
        > https://ginoauth-example.herokuapp.com/callback/google?state=xxxx&code=@#!@%!%!
        <br>

        跳轉到callback，同時會將state跟authorization code加入跳轉url的query string之中。  
        其中state參數必須跟第一步驟跳轉到google時填入的state參數相同，避免CSRF攻擊。  
        通常state會存在session內，以便於比對。
        
        > 詳情請參考[OAuth 2.0 筆記 (7) 安全性問題](https://blog.yorkxin.org/2013/09/30/oauth2-7-security-considerations.html)
        <br>

        #### 藉由authorization code取得access code
        再來談談code參數，後端跳轉到callback並且接受code之後，需要再以code參數去跟google拿Token。  
        再藉由Token去跟google獲得真正想要的資源。

        > 向google申請Access Token的url<br>
        > https://oauth2.googleapis.com/token?client_id=xxx
        > &client_secret=xxxx  
        > &authorization_code=xxxxx  
        > &grant_type=code

        參數說明
        > - client_id, client_secret: 讓google得知此次請求是由哪個第三方應用發送。
        > - authorization_code: 剛剛上方取得的授權碼。
        > - grant_type: 告知此次請求是採用Authorization-Code的方式進行。
        
        回傳json形式範例(根據要求資源的公司不同，回傳的json格式會有些許落差。)
        ```json
        {
          "access_token": "xxxx", 
          "token_type": "xxx",
          "expires_in": "放token過期時間",
          "refresh_token": "當access token過期時，可以再次申請一組Access Token",
          "scope": "權限",
          "uid": "xxxx",
          "info": "xxxx",
        }
        ```

        最後就可以藉由上述流程中獲得的access_token來申請資源啦。  
        要注意的是**client_secret, access_token等資訊一定要放在後端並且小心存放避免暴露**。

    - ### Implicit(隱藏式):
      
      這個方法現在已經比較少用了，應用場景主要是在沒有後端的純前端應用上。

      直接將Access Token儲存在前端， 並且因為隱藏了授權碼的使用， 所以被稱為隱藏式。

      由於此方法非常的不安全， 會將Access Token暴露在前端， 所以通常只能用在安全需求低的場景， 而且對於Token的有效時間要設置的非常短， 避免有心人士利用。

      通過這個方式得到授權之後，google會根據申請時候的redirect參數去跳轉。 範例如下
      > https://test.com/callback#token=Access_Token

      值得注意的是回傳的Access_Token並不是放在Query_String中，而是放在了Fragment內， 這樣的好處是參數並不會傳回伺服器，而是僅供瀏覽器使用， 可以防止中間人攻擊的發生。

      更多有關URL參數以及Fragment的細節可以參考: [URL](https://terrylee7788.wordpress.com/2015/07/11/url-%E7%9A%84%E5%90%84%E5%80%8B%E9%83%A8%E5%88%86/) [Wiki](https://en.wikipedia.org/wiki/URI_fragment)

    - ### Password(密碼式)

      這種情境個人認為又更少使用了， 此認證方式只會出現在當使用者今天高度的信任某個第三方應用程式， RFC 6749也允許在此情況下直接將用戶google的帳號密碼告訴第三方的應用程式。

      使用這種方式的時候Access_Token會直接以json的方式回傳給第三方的應用，此種方式通常較不推薦， 畢竟在網路的世界，還是要以安全為最高準則。 

    - ### Client Credentials(憑證式)

      最後一種認證方式我們稱之為憑證式，運用的場景就是今天的第三方應用並不存在前後端，而是在command底運行的應用程式，這種認證方式通常不是針對使用者，而是針對第三方應用程式量身打造的。

      一般像是用在，例如我今天寫了一個應用程式需要接上youtube的API，這樣的話我就需要只用這種認證方式，讓youtube知道今天我的某個應用程式要存取youtube的資源。

      我寫的應用程式存取到youtube的資源之後再給全部的使用者使用， 有點類似底下的使用者都可以共享我的應用程式所存取的資源。

- ## 第三方登入流程簡述
 
  - 今天假設我們的第三方登入網站為A， 開放用戶藉由google帳號去做登入，流程大概如下:

  > A跳轉至google登入界面。<br>
  > ↓<br>
  > google要求用戶登入google帳號，並且詢問是否願意開放權限給A做使用。<br>
  > ↓<br>
  > 用戶同意開放權限，google再跳轉回A網站，並且在query string附上code。<br>
  > ↓<br>
  > A網站在後端使用code向google申請Access Token。<br>
  > ↓<br>
  > google將Access Token以及一些相關資訊放入json回傳。<br>
  > ↓<br>
  > A網站藉由Access Token向google請求相關的用戶資源。
  
- ## Reference
  - [GitHub OAuth 第三方登入範例](http://www.ruanyifeng.com/blog/2019/04/github-oauth.html)
  - [Google Sign-in with OAuth 2.0](https://yushuanhsieh.github.io/post/2018-08-25-go-google-oauth/#comments)
  - [OAuth 2.0筆記系列](https://blog.yorkxin.org/2013/09/30/oauth2-4-1-auth-code-grant-flow.html)
  - [google oauth2](https://developers.google.com/identity/protocols/oauth2/web-server#httprest_7)
  
  