import React from 'react';
import styled from 'styled-components';
import Header from './Header';

export default function App() {
  const App = styled.div`
    box-sizing: border-box;
    background-color: rgb(21, 21, 21);
    & *{
      box-sizing: border-box;
    }
  `;
  const Top = styled.div`
    border-bottom: 1px solid rgb(52, 52, 52);
    height: 48px;
  `;
  const Body = styled.div`
    display: flex;
  `;
  const LeftPanel = styled.div`
    height: 48px;
    width: 100px;
    border-right: 1px solid rgb(52, 52, 52);
  `;
  return (
    <App>
      <Top><Header/></Top>
      <Body>
        <LeftPanel/>
        <div/>
      </Body>
    </App>
  );
}
